package parser_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/parser"
	mock_parser "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/parser/mockParser"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewsQuerier(t *testing.T) {
	_, err := http.Get("https://example.com")

	if err == nil {
		querier, err := parser.NewQuerier(
			parser.WithDefaultHeader(),
			parser.WithDefaultClient(),
			parser.WithClientDefaultTimeout(),
		)
		require.NoError(t, err)
		require.NotNil(t, querier)
		require.True(t, querier.HasContentEncodingHandler("gzip"))

		rawURLs := []string{
			"https://news.pts.org.tw/article/652813",
			"https://www.cna.com.tw/news/aspt/202309050416.aspx",
		}

		for i := range rawURLs {
			rawURL := rawURLs[i]
			t.Run(rawURL, func(t *testing.T) {
				q := querier.DoQuery(parser.NewQuery(rawURL))
				require.NoError(t, q.Error)

				content, err := q.Content()
				require.NoError(t, err)

				var body []byte
				body, q.Error = io.ReadAll(content)
				require.NoError(t, err)
				require.NotNil(t, body)

				v := validator.New()
				require.NoError(t, v.Var(string(body), "html"))
			})
		}
	} else {
		t.Log("Newwork is unreachable, skip this test")
	}
}

func TestQueryPipeline(t *testing.T) {
	htmlBody := `<!DOCTYPE html>
       <html lang="en">
       <head>
               <meta charset="UTF-8">
               <meta name="viewport" content="width=device-width, initial-scale=1.0">
               <title>Document</title>
       </head>
       <body>
               <h1>For testing</h1>
       </body>
	    </html>`

	type testCase struct {
		Id         int
		Path       string
		StatusCode int
		TimeDelay  time.Duration
	}

	tcs := []testCase{
		{
			Path:       "/ok",
			StatusCode: http.StatusOK,
			TimeDelay:  0,
		},
		{
			Path:       "/notfound",
			StatusCode: http.StatusNotFound,
			TimeDelay:  0,
		},
		{
			Path:       "/timeout",
			StatusCode: http.StatusOK,
			TimeDelay:  2 * time.Second,
		},
		{
			Path:       "/unauthorized",
			StatusCode: http.StatusUnauthorized,
			TimeDelay:  0,
		},
	}

	path2id := map[string]int{}

	mux := chi.NewRouter()
	mux.Use(middleware.Compress(5, "gzip"))

	HandlerFromTC := func(tc testCase) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			t.Logf("hit: %s Handler wait for %s sec\n", tc.Path, tc.TimeDelay)
			time.Sleep(tc.TimeDelay)
			w.WriteHeader(tc.StatusCode)
			if tc.StatusCode == http.StatusOK {
				w.Write([]byte(htmlBody))
			}
			return
		}
	}

	for i := range tcs {
		tcs[i].Id = i
		path2id[tcs[i].Path] = tcs[i].Id
		mux.Get(tcs[i].Path, HandlerFromTC(tcs[i]))
	}
	srvc := httptest.NewServer(mux)

	inChan := make(chan *parser.Query)
	querier, err := parser.NewQuerier(
		parser.WithDefaultHeader(),
		parser.WithDefaultClient(),
		parser.WithClientTimeout(1*time.Second),
	)
	require.NoError(t, err)
	require.NotNil(t, querier)

	ctx := context.Background()
	outChan, errChan := querier.DoQueryPipeline(ctx, inChan)

	go func() {
		defer close(inChan)
		for i := range tcs {
			inChan <- parser.NewQueryWithId(tcs[i].Id, srvc.URL+tcs[i].Path)
		}
	}()

	v := validator.New()
	for i := 0; i < len(tcs); i++ {
		select {
		case q := <-outChan:
			require.NoError(t, q.Error)
			require.Equal(t, path2id["/ok"], q.Id())
			require.Equal(t, 200, q.RespHttpStatusCode())

			rc, err := q.Content()
			require.NoError(t, err)
			require.NotNil(t, rc)

			doc, err := io.ReadAll(rc)
			require.NoError(t, err)
			require.NotNil(t, doc)

			require.NoError(t, v.Var(string(doc), "html"))
		case err := <-errChan:
			require.Error(t, err)
			es := err.Error()

			if strings.HasPrefix(es, fmt.Sprintf("%d-th", path2id["/timeout"])) {
				require.True(t, strings.Contains(es, "context deadline exceeded"))
			}

			if strings.HasPrefix(es, fmt.Sprintf("%d-th", path2id["/notfound"])) {
				require.True(t, strings.Contains(es, "request error with error code 404"))
			}

			if strings.HasPrefix(es, fmt.Sprintf("%d-th", path2id["/unauthorized"])) {
				require.True(t, strings.Contains(es, "request error with error code 401"))
			}
		}
	}
}

func TestUDNParser(t *testing.T) {
	p := parser.NewUDNParser()

	type testCase struct {
		FileName string
		Title    string
		Category string
		Author   []string
		GUID     string
		Tags     []string
	}

	tcs := []testCase{
		{
			FileName: "news/001.html",
			Title:    "核汙水排入海 中日緊張升溫",
			Category: "全球",
			Author:   []string{"茅毅", "謝守真"},
			GUID:     "123707-7398504",
			Tags:     []string{"福島", "核汙水", "日本"},
		},
		{
			FileName: "news/002.html",
			Title:    "家暴說是否道歉？楊志良晚間發表聲明 對這群人道歉了",
			Category: "要聞",
			Author:   []string{"沈能元"},
			GUID:     "6656-7400449",
			Tags:     []string{"2024選舉", "2024總統選舉", "楊志良", "郭台銘"},
		},
		{
			FileName: "news/003.html",
			Title:    "新聞幕後／侯友宜與郭台銘前天見過面 預計周三咖啡會",
			Category: "要聞",
			Author:   []string{"張睿廷"},
			GUID:     "123307-7400824",
			Tags:     []string{"2024選舉", "2024總統選舉"},
		},
		{
			FileName: "global/001.html",
			Title:    "當選的代價？泰國新總理斯雷塔，為泰黨向保皇派低頭的執政賭局",
			Category: "政經角力",
			Author:   []string{"徐子軒"},
			GUID:     "8663-7394459",
			Tags:     []string{"泰國", "深度專欄"},
		},
		{
			FileName: "global/002.html",
			Title:    "北京地產暴雷危機？SOHO中國欠稅案，與「潤」到美國的創辦人潘石屹",
			Category: "過去24小時",
			Author:   []string{"轉角24小時"},
			GUID:     "8662-7389296",
			Tags:     []string{"過去24小時", "中國", "美國", "經濟"},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.FileName,
			func(t *testing.T) {
				t.Parallel()
				f, err := os.Open(fmt.Sprintf("example/udn/%s", tc.FileName))
				require.NoError(t, err)
				q := p.Parse(parser.NewTestQuery(200, "", f))
				require.NoError(t, q.Error)
				require.NotNil(t, q.News)

				require.Equal(t, tc.GUID, q.News.GUID)
				require.Equal(t, tc.Title, q.News.Title)
				require.Equal(t, tc.Category, q.News.Category)
				require.False(t, q.News.PubDate.IsZero())
				require.ElementsMatch(t, tc.Author, q.News.Author)
				require.ElementsMatch(t, tc.Tags, q.News.Tag)
			},
		)
	}
}

func TestCTParser(t *testing.T) {
	p := parser.NewCTParser()
	type testCase struct {
		FileName string
		Title    string
		Category string
		Author   []string
		GUID     string
		Tags     []string
		Content  []string
	}

	tcs := []testCase{
		{
			FileName: "001.html",
			Title:    "網紅行動電源爆炸 舊換新影響25萬人 業者緊急公告再道歉",
			Category: "科技",
			Author:   []string{"陳駿碩"},
			GUID:     "20230825004474-260412",
			Tags:     []string{"墨子科技", "行動電源", "充電器", "網紅行動電源", "網紅"},
			Content: []string{
				"有許多網紅推薦、團購的行動電源日前發生爆炸意外，廠商墨子科技今（25）日發文公布兩組產品序號讓用戶可以免費更換，而今晚業者又發布緊急公告，表示有高達25萬名用戶受影響，目前線上客服及電話訊息量已滿載，請用戶耐心等待。",
				"「MOZTECH墨子科技」今日晚間在臉書粉專發布的緊急公告指出，此次特殊的批次換新服務約有25萬名用戶受到影響，因此公司的線上客服、電話訊息皆已滿載，對於難以短時間內提供完善的售後服務，墨子科技再次道歉。",
				"墨子科技表示，目前已經在全台新增25間授權獨立維修中心，以及MOZTECH墨子科技三創體驗店，讓用戶可以批次換新，該服務將於8月30日啟用，使用者也可於當天前往更換。",
				"另外，特殊色「奶油黃」為聯名限量色，尚需等待45到50天的生產工作天；而特殊色「奶茶色」則是因為原物料已停產，故無法另行生產，將以燕麥奶替代換新。",
				"事實上，墨子科技今日上午就發布公告，由於目前尚未得知產品檢測結果、產品批次號碼，為避用戶的使用疑慮，若消費者持有產品規格序列號202211、202212的「MOZTECH® 萬能充Pro 多功能五合一行動電源」，發現產品有任何瑕疵問題，經檢測單位確認後，墨子科技將直接免費提供一台全新的產品更換給用戶。",
			},
		},
		{
			FileName: "002.html",
			Title:    "星巴克連兩天買一送一！2款星冰樂入列",
			Category: "生活",
			Author:   []string{"莊楚雯"},
			GUID:     "20230828001159-260405",
			Tags:     []string{"星巴克", "買一送一", "星冰樂", "夏天", "伯朗咖啡"},
			Content: []string{
				"颱風還沒來，今、明兩天仍是高溫炎熱的夏季天氣，星巴克從今、明2天推出特定品項飲料買一送一優惠，品項包含那堤、夏日冰柚冷萃咖啡、焦糖咖啡星冰樂、醇濃抹茶奶霜星冰樂等，活動時間為上午11點晚上8點。",
				"星巴克表示，今、明（28、29日）2天推出指定品項好友分享日，於上午11點至晚上8點，至門市購買2杯特大杯冰熱／風味一致的指定飲料，享買一送一優惠，活動品項為那堤、特選那堤、焦糖瑪奇朵、特選焦糖瑪奇朵、冷萃咖啡、夏日冰柚冷萃咖啡、椰奶經典巧克力、豆奶玫瑰蜜香茶那堤、焦糖咖啡星冰樂、醇濃抹茶奶霜星冰樂。不包含罐裝飲料、典藏系列咖啡、手沖、虹吸式咖啡及含酒精飲料。",
				"星巴克表示，買一送一活動不適用於車道服務、外送外賣、電話預訂與行動預點服務，也不適用機場門市、高鐵門市、松山車站、台鐵一、板橋車站、中壢休息站、湖口休息站、西湖休息站、泰安南、泰安北、南投休息站、西螺、東山休息站、關廟南、仁德南、仁德北、草山、花蓮和平、洄瀾、宜蘭頭城、福隆觀海、清境、日月潭、中埔穀倉、奇美博物館、墾丁福華、海生館、龍門、101典藏。",
				"此外，路易莎周三（30日）前也有好友分享日，全門市每周一、二、三上午11點前全飲料第二件半價，不含手沖及精品咖啡；伯朗咖啡卡友在今天也有買一送一，持伯朗Lounge卡消費，可享飲品買一送一，飲品不限品項與容量，結帳金額以高價品計算，手沖咖啡、瓶裝飲品、啤酒不適用優惠活動。",
			},
		},
		{
			FileName: "003.html",
			Title:    "《華爾街日報》：習近平反對西式消費刺激經濟增長模式",
			Category: "兩岸",
			Author:   []string{"陳柏廷"},
			GUID:     "20230828002197-260409",
			Tags:     []string{"華爾街日報", "習近平", "西式消費", "刺激", "經濟增長"},
		},
		{
			FileName: "004.html",
			Title:    "42架F-16戰機給烏克蘭 它要討回9年前慘劇公道",
			Category: "軍事",
			Author:   []string{"楊幼蘭"},
			GUID:     "20230828001313-260417",
			Tags:     []string{"F-16", "戰機", "俄羅斯", "普丁", "烏克蘭"},
		},
		{
			FileName: "005.html",
			Title:    "畢業即失業？ 他求職四處碰壁 青年計畫協助進入職場",
			Category: "生活",
			Author:   []string{"林欣儀"},
			GUID:     "20230828003954-260405",
			Tags:     []string{"就業", "勞動部", "培訓", "職場", "青年"},
		},
		{
			FileName: "006.html",
			Title:    "新聞透視》黑道變網紅 斷金流才是反制之道",
			Category: "社會新聞",
			Author:   []string{"胡欣男"},
			GUID:     "20230828000364-260106",
			Tags:     []string{"黑道", "掃黑", "一清專案", "明仁", "黑幫"},
			Content: []string{
				"40年前一清專案，意外造就黑金政治，如今黑道不僅愈掃愈多，還搞行銷做網紅，流量比警察預防犯罪宣導，高不知幾百倍，政府面子該往哪擺？黑道網紅化，追根究柢是有黑金護航，黑白掛勾弊案叢生，黑道毫不避諱行銷，掃黑怎會有效？",
				"部分黑道大哥一清專案出獄後，從政漂白，招來黑道治國、黑金政治的批評。去年底九合一大選，民進黨陷入「黑道中常委」爭議，接連發生台北88會館、台南88槍擊案，老百姓對黑金反感，成為綠營敗選主因之一。",
				"為此，民進黨主席賴清德不惜得罪部分派系人士，也要修選罷法排黑，為明年大選拆彈。然而不論法怎麼修，始終是針對參選者，或許杜絕黑道參政，但盤根錯節的黑金政治結構，多的是隱身幕後的影武者，沉痾由來已久，豈是幾條法令能根除。",
				"黑金問題未解，再有竹聯幫明仁會的奢華春酒，透過網路瘋傳，政府高層才驚覺，新生代黑道搞詐欺後的高調炫富，根本無視公權力。而後雖修訂「明仁會條款」，黑道在婚喪喜慶公開活動「表現空間有限」，卻仍能在網路發展。",
				"網路短影音隨處可見黑道題材，不僅賺流量，更紊亂社會價值觀，彷如黑道招生，對岸都驚嘆台灣黑道的高調。當警方移送畫面被當作黑幫抬高江湖身價的素材，叫辛苦掃黑的警察情何以堪。",
				"民主國家不太可能管制人民的網路行為，但黑幫高調宣揚，就是挑戰公權力。司法實務上，黑道要真的被重判入獄，屈指可數，唯有透過科技偵查斷其金流，讓黑道沒錢囂張，既怕又痛，才能反制黑幫氣焰。",
			},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.FileName,
			func(t *testing.T) {
				t.Parallel()
				f, err := os.Open(fmt.Sprintf("example/ct/%s", tc.FileName))
				require.NoError(t, err)
				q := p.Parse(parser.NewTestQuery(200, "", f))
				require.NoError(t, q.Error)
				require.NotNil(t, q)
				require.Equal(t, tc.GUID, q.News.GUID)
				require.Equal(t, tc.Title, q.News.Title)
				require.Equal(t, tc.Category, q.News.Category)
				require.False(t, q.News.PubDate.IsZero())
				require.ElementsMatch(t, tc.Author, q.News.Author)
				require.ElementsMatch(t, tc.Tags, q.News.Tag)
				if len(tc.Content) > 0 {
					require.Equal(t, len(tc.Content), len(q.News.Content))
					for i, c := range q.News.Content {
						require.Equal(t, tc.Content[i], c)
					}
				}
			},
		)
	}
}

func TestPTSParser(t *testing.T) {
	p := parser.NewPTSParser()

	type testCase struct {
		FileName string
		Title    string
		Category string
		Author   []string
		GUID     string
		Tags     []string
		Content  []string
	}

	tcs := []testCase{
		{
			FileName: "001.html",
			Title:    "美韓軍演首納太空軍 一文了解太空軍是什麼",
			Category: "全球",
			Author:   []string{"陳宥蓁"},
			GUID:     "652596",
			Tags:     []string{"空軍", "美國", "衛星", "部隊", "駐韓美軍"},
			Content: []string{
				"美國和南韓於昨（21）日起為期10天舉行「乙支自由之盾」年度聯合軍演，為防堵北韓網路入侵，還首度將美軍成立不久的「太空軍」加入演習。駐韓美軍先前指出，太空軍的加入會使作戰更具多樣性。",
				"美國太空軍是什麼？",
				"美國軍隊共分為6大軍種，太空軍於2019年12月成立，屬最新分支，總編制約為1萬6千人，主要工作是組織、訓練和裝備太空部隊，不過在先期任務規劃還未包括直接作戰，為的是保護美國在太空的資產及利益，阻止他人藉由太空進行侵略。",
				"在太空軍成立之際，美國前總統川普曾說，「太空是世界上最新的戰鬥領域。」而今年年初，太空軍作戰部長薩茲曼便指出，中國和俄羅斯已發展出反衛星飛彈和軌道攔截能力等技術，是美國在太空領域的兩大威脅。",
				"美國太空軍如何分工？",
				"太空軍作戰部長辦公室（OCSO）是太空軍的最高指揮機構，其下有太空軍戰地司令部（Space Force Field Command），3個司令部又有各自的三角洲部隊（Delta），並有中隊（Squadron）提供支援。",
				"太空作戰司令部（SpOC） 負責做好戰鬥準備，監視與偵察全球情報，並與外界夥伴合作。",
				"太空系統司令部（SSC） 負責為開發、購買和部署太空武裝系統，監督國防部衛星和其他空間系統的發射操作、維護等，與太空相關的研究工作。",
				"太空訓練和準備司令部（STARCOM） 負責培訓和教育人員，並制定太空術語和戰術等。",
				"新單位用「死神」當隊徽？",
				"美國太空軍近期成立了新單位ISRS，是首支恐會摧毀敵方衛星的部隊，與採用「死神」當隊徽的意涵相互呼應；而死神眼中的北極星芒，則代表安全的導引。",
				"ISRS負責的任務是阻斷對美國使用衛星系統造成的威脅，包含地面上可能造成影響的雷射、干擾訊號裝置，或駭入衛星系統的資安攻擊等。他們將分析潛在目標，並進行識別和追蹤。",
				"美國境外也有太空軍？",
				"美國太空軍去（2022）年11月在夏威夷珍珠港的基地，首度在本土以外設立區域指揮中心，除了控管印太地區的太空軍事行動，也是為了防禦中國威脅，並包括俄羅斯及北韓。",
				"而同年12月，駐韓美軍也成立了美國太空軍的海外太空部隊，負責監控和追蹤北韓及鄰近地區飛彈，還有全球定位、衛星通訊等任務。",
			},
		},
		{
			FileName: "002.html",
			Title:    "雲林垃圾轉換衍生燃料堆置 縣府盼修法放寬添加率",
			Category: "環境",
			Author:   []string{"王威雄"},
			GUID:     "652936",
			Tags:     []string{"垃圾", "燃料棒", "環保局", "雲林縣"},
			Content: []string{
				"將垃圾送進零廢棄資源化系統，轉換成衍生燃料俗稱燃料棒，能取代過去的生煤作為鍋爐摻配使用，為了解決雲林縣垃圾危機，雲縣府2020年引進這套系統，但現在也碰上燃料棒去化的問題。",
				"雲林縣環保局長張喬維表示，「目前雲林縣一天大約有150噸的SRF（衍生燃料）產生量，現在已堆置到1萬多噸，主要是塑化公司只有一套設備來進行操作。」",
				"雲林縣環保局表示，雲林縣將衍生燃料棒以有價方式販售給廠商，但因價格高市場銷售受限，連帶也影響到以零廢棄資源化系統無法全力運轉，希望中央能針對運用衍生燃料再生能源，訂定相關獎勵辦法。",
				"張喬維提及，「六輕這邊混燒達5%，要再添加到10%、甚至20%，有一些法令上的限制，所以目前我們也請中央針對添加到20%來進行修法。」",
				"環保局表示，中央訂出2050年淨零排放的目標，但目前未針對減碳效益提出進度，許多廠商保持觀望態度，不願意投資改善設備，預估年底六輕塑化將再新增機組去化燃料棒，希望中央機關能協助相關法令規範，推廣廢棄物資源化。",
			},
		},
		{
			FileName: "003.html",
			Title:    "印度月船3號創全球登月球南極之先 尋水冰與地質分析",
			Category: "全球",
			Author:   []string{"鄭惟仁", "鍾建剛"},
			GUID:     "652902",
			Tags:     []string{"印度", "探測器", "月球", "登月", "登陸"},
		},
		{
			FileName: "004.html",
			Title:    "新北酒駕男撞老翁致重傷搶救 車上3人遭連坐罰",
			Category: "社會",
			Author:   []string{"陳冠勳", "張梓嘉", "溫正衡"},
			GUID:     "652946",
			Tags:     []string{"台南地檢署", "過馬路", "酒駕"},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.FileName,
			func(t *testing.T) {
				// t.Parallel()
				f, err := os.Open(fmt.Sprintf("example/pts/%s", tc.FileName))
				require.NoError(t, err)
				q := p.Parse(parser.NewTestQuery(200, "", f))
				require.NoError(t, q.Error)
				require.NotNil(t, q)
				require.Equal(t, tc.GUID, q.News.GUID)
				require.Equal(t, tc.Title, q.News.Title)
				require.Equal(t, tc.Category, q.News.Category)
				require.False(t, q.News.PubDate.IsZero())
				require.ElementsMatch(t, tc.Author, q.News.Author)
				require.ElementsMatch(t, tc.Tags, q.News.Tag)

				if len(tc.Content) > 0 {
					require.Equal(t, len(tc.Content), len(q.News.Content))
					for i, c := range q.News.Content {
						require.Equal(t, tc.Content[i], c)
					}
				}
			},
		)
	}
}

func TestCNAParser(t *testing.T) {
	p := parser.NewCNAParser()

	type testCase struct {
		FileName string
		Title    string
		Category string
		Author   []string
		GUID     string
		Tags     []string
		Content  []string
	}

	tcs := []testCase{
		{
			FileName: "001.html",
			Title:    "新北攤商涉用豬肉混充羊肉 最重可罰400萬加坐牢",
			Category: "生活",
			Author:   []string{"李亨山"},
			GUID:     "ahel-202308230179",
			Tags:     []string{"林金富", "食藥署", "豬肉", "劉芳銘"},
			Content: []string{
				"（中央社記者沈佩瑤台北23日電）媒體報導，位於新北市的大台北果菜市場某攤商販售羊肉遠低於市價，涉以豬肉魚目混珠。食藥署今天表示，已前往稽查並抽驗，若確實摻偽假冒，最高可罰400萬元並面臨坐牢。",
				"根據CTWANT網站今天報導，新北市蘆洲區的大台北果菜市場某攤商販售羊肉價格遠低於市價，二度送台灣檢驗科技股份有限公司（SGS）檢驗後，發現羊肉根本不含任何羊肉成分，而是用豬肉混充。",
				"衛生福利部食品藥物管理署北區管理中心主任劉芳銘在例行記者會中說明，新北市衛生局昨天接獲消息已立即前往稽查，了解現場及上游供應情形，現場原料包括羊肉、豬肉、牛肉，羊肉來自紐西蘭及澳洲。",
				"是否有以豬肉混充羊肉，劉芳銘表示，現場確實有販售羊肉，已抽驗半成品、羊肉火鍋肉片送檢，檢驗結果預計5到7天出爐。",
				"食藥署副署長林金富指出，若送驗後發現確實有摻偽假冒情形，除可依食品安全衛生管理法，針對標示不實開罰新台幣4萬到400萬元，另外也可移送檢調偵辦，因涉違反刑法255條意圖欺詐他人等情事。",
				"過去國內也曾發生羊肉混充案件，林金富表示，民國109年以來共抽驗55件，發現3件羊肉含有豬肉成分，其中2件來自同樣源頭，地方衛生局皆依法裁罰，總計開罰8萬元。",
				"根據中華民國刑法255條規定，意圖欺騙他人，而就商品之原產國或品質，為虛偽之標記或其他表示者，處1年以下有期徒刑、拘役或3萬元以下罰金。（編輯：李亨山）1120823",
			},
		},
		{
			FileName: "002.html",
			Title:    "中國單身人口逾2億 專家說因年輕人太宅",
			Category: "兩岸",
			Author:   []string{"唐佩君", "呂佳蓉"},
			GUID:     "acn-202308230192",
			Tags:     []string{"中國", "環球時報"},
			Content: []string{
				"（中央社台北23日電）中國年輕未婚單身人口達2.39億人，且初婚年齡延後至28.67歲；有專家表示，調查年輕人為何不婚不戀，原因包括社交圈子固定、太宅及不喜社交，但有網友直言，原因很簡單就是窮。",
				"8月22日是農曆七夕情人節，綜合陸媒報導，據2022年中國統計年鑑數據顯示，截至2021年，全國15歲以上單身人口約為2.39億人。",
				"此外，中國年輕人婚育年齡也普遍推遲，據2020年「中國人口普查年鑑」顯示，平均初婚年齡28.67歲，比2010年的24.89歲增加了3.78歲。",
				"中國單身人群越來越龐大，中國科學院心理研究所簽約心理諮詢師鄭莉接受環球時報採訪時說，根據調查，年輕人不婚不戀的理由，主要是社交圈子固定、太宅不喜社交、不善表達自己及閒暇時間忙於上網等。",
				"鄭莉認為，長久以來，中國教育制度以智力學歷教育為主，從小到大學培養的都是各種生存的能力，而與戀愛、情商等相關的知識灌輸極少，造成年輕人普遍缺乏談情說愛的能力，往往面臨上學期間父母禁止戀愛，畢業之後父母催促結婚的尷尬局面。",
				"她認為，解決年輕人婚戀難的問題，單身青年要從根本上著手，即主動接受成人情感教育，補上孩童時代所缺的情感課，並規劃好社交表格，強迫自己打破固定的社交圈及訓練表達能力。",
				"不過，部份網友不認同這番評論，網路上不少留言表示，最直接原因就是「窮」，「生活太難」、「結不起」、「沒錢」等。",
				"也有部份網友質疑這項統計，指15歲才高一就納入未婚人口統計範圍是不是太廣了，才導致單身人口如此高。（編輯：唐佩君/呂佳蓉）1120823",
			},
		},
		{
			FileName: "003.html",
			Title:    "庫克群島總理：福島核處理水排放 太平洋島國存歧見",
			Category: "國際",
			Author:   []string{"王嘉語", "嚴思祺"},
			GUID:     "aopl-202308230260",
			Tags:     []string{"核處理水", "IAEA", "國際原能總署", "福島", "布朗"},
			Content: []string{
				"（中央社雪梨23日綜合外電報導）庫克群島總理、「太平洋島國論壇」輪值主席布朗表示，日本將來自福島核電廠核處理水排放入海的決定獲科學支持，但太平洋地區可能不會對這個「複雜」議題達成共識。",
				"路透社報導，日本政府昨天表示，將於24日開始將超過100萬公噸核處理水排入海中。這些核處理水來自被毀的福島第一核電廠，排放計畫受到中國的嚴厲批評。",
				"日本表示，核處理水排放是安全的。聯合國核監督機構國際原子能總署（International Atomic Energy Agency, IAEA）今年7月為排放計畫背書，表示這項排放符合國際標準，對人類和環境的影響「微不足道」。",
				"國際原子能總署於7月前往庫克群島（Cook Islands），向「太平洋島國論壇」（Pacific Islands Forum）提交相關調查結果。",
				"「太平洋島國論壇」是由18個區域國家組成的集團，其專屬經濟區橫跨太平洋，面積達到4000萬平方公里，全球有一半的鮪魚漁獲來自這個地區。",
				"布朗（Mark Brown）今天在聲明中表示：「我相信排放符合國際安全標準。」他並說，國際原子能總署將在排放過程中繼續監控核處理水。",
				"但布朗也提到，並非所有太平洋地區領導人都持相同立場，「太平洋島國論壇」可能不會達成一致的立場。",
				"布朗表示，在曾受到外部勢力進行核武器試驗所影響的太平洋地區，這是一項「複雜的議題」。美國在1940年代和1950年代間，法國在1966年至1996年間，於太平洋島國進行了核試驗。",
				"太平洋無核區在1985年根據太平洋非核區條約（Pacific Nuclear Free Zone Treaty）建立，防止放射性物質的傾倒。",
				"斐濟總理拉布卡（Sitiveni Rabuka）21日在談話中表示，基於國際原子能總署的報告，他支持排放，並認為將持續30年受管控的核處理水排放，與在太平洋地區進行的核武器實驗連結起來，是在「散佈恐懼」。",
				"美拉尼西亞先鋒集團（Melanesian Spearhead Group）將於明天舉行會議，討論福島核處理水排放問題。這個集團的成員包含巴布亞紐幾內亞、斐濟、萬那杜、索羅門群島和新喀里多尼亞島執政黨「卡納克社會主義民族解放陣線」（FLNKS）。（譯者：王嘉語/核稿：嚴思祺）1120823",
			},
		},
		{
			FileName: "004.html",
			Title:    "泰王任命賽塔為總理 新政府親軍方成員備受爭議",
			Category: "國際",
			Author:   []string{"何宏儒", "楊昭彥"},
			GUID:     "aopl-202308230370",
			Tags:     []string{"帕拉育", "政變", "泰國", "戴克辛", "為泰黨"},
		},
		{
			FileName: "005.html",
			Title:    "數位部：推動數據公益恪遵個資法 確保當事人權益",
			Category: "科技",
			Author:   []string{"楊凱翔"},
			GUID:     "ait-202308230354",
			Tags:     []string{"台權會", "數位部", "個資法", "歐盟"},
		},
		{
			FileName: "006.html",
			Title:    "高雄女監收容人突昏迷不治 獄方否認做登革熱防疫過勞死",
			Category: "社會",
			Author:   []string{"陳仁華"},
			GUID:     "asoc-202308230348",
			Tags:     []string{"防疫", "高雄女子監獄", "登革熱", "高雄市"},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.FileName,
			func(t *testing.T) {
				// t.Parallel()
				f, err := os.Open(fmt.Sprintf("example/cna/%s", tc.FileName))
				require.NoError(t, err)
				q := p.Parse(parser.NewTestQuery(200, "", f))
				require.NoError(t, q.Error)
				require.NotNil(t, q)
				require.Equal(t, tc.GUID, q.News.GUID)
				require.Equal(t, tc.Title, q.News.Title)
				require.Equal(t, tc.Category, q.News.Category)
				require.False(t, q.News.PubDate.IsZero())
				require.ElementsMatch(t, tc.Author, q.News.Author)
				require.ElementsMatch(t, tc.Tags, q.News.Tag)
				require.NotContains(t, q.News.Tag, "NewsArticle")
				if len(tc.Content) > 0 {
					require.Equal(t, len(tc.Content), len(q.News.Content))
					for i, c := range q.News.Content {
						require.Equal(t, tc.Content[i], c)
					}
				}
			},
		)
	}
}

func TestLTNParser(t *testing.T) {
	p := parser.NewLTNParser()

	type testCase struct {
		FileName string
		Title    string
		Category string
		Author   []string
		GUID     string
		Tags     []string
		Content  []string
	}

	tcs := []testCase{
		{
			FileName: "001.html",
			Title:    "瓦格納首腦疑墜機亡 航跡追蹤網站：墜機前30秒出現問題",
			Category: "國際",
			Author:   []string{"即時新聞"},
			GUID:     "news-world-breakingnews-4405752",
			Tags:     []string{"萊格賽600", "墜機問題", "莫斯科", "瓦格納", "普里格津", "Flightradar24", "瓦格納集團", "瓦格納傭兵", "地對空飛彈", "飛機墜毀", "特維爾", "航跡追蹤網站", "俄羅斯", "瓦格納首腦"},
			Content: []string{
				"〔即時新聞／綜合報導〕1架巴西航空工業公司製造的「萊格賽600」（Embraer Legacy 600）私人飛機，23日在莫斯科北部的特維爾（Tver）地區墜毀，機上10人全數罹難，據傳機上載有俄國傭兵組織瓦格納集團首腦普里格津（Yevgeny Prigozhin）。航跡追蹤網站「flightradar24」指出，此架飛機在墜毀前30秒以前都很正常，事發前30秒，飛機突然垂直下墜。",
				"根據《路透》報導，「flightradar24」主編佩琴尼克（Ian Petchenik）指出，在當地時間傍晚6點19分左右，飛機突然垂直向下墜落，大約在30秒內，飛機從2.8萬英呎高度驟降8000多英呎，「無論發生了些什麼，這一切都來的很快，但在飛機急劇下降之前，沒有任何跡象顯示這架飛機有問題」。",
				"「flightradar24」指出，這架飛機在30秒內經歷一系列高度上升或下降的劇烈操作，最後發生災難性後果，「flightradar24」最後於俄羅斯當地時間傍晚6點20分接收到飛機的最後數據。",
				"部分消息人士告訴俄羅斯媒體，他們認為這架飛機是被1枚或多枚地對空飛彈擊落。《路透》表示，對此說法無法查證。一位知情人士指稱，他在「Flightradar24」上發現這架私人飛機的註冊號碼為RA-02795，與兵變後載著普里格津飛往白俄羅斯的飛機相同。",
			},
		},
		{
			FileName: "002.html",
			Title:    "日本核廢水正式排海 原能會：啟動放射性物質擴散預報系統",
			Category: "生活",
			Author:   []string{"吳柏軒"},
			GUID:     "news-life-breakingnews-4406757",
			Tags:     []string{"日本", "核廢水", "日本核廢水"},
		},
		{
			FileName: "003.html",
			Title:    "中職》隊友全楞住！江坤宇驚喜再見轟 兄弟休息室手忙腳亂網笑翻",
			Category: "中職",
			Author:   []string{"體育中心"},
			GUID:     "news-breakingnews-4405889",
			Tags:     []string{"李振昌", "延長賽", "郭明錤", "江坤宇", "張家齊", "兄弟休息室", "中職", "馬拉松", "頻道", "再見轟", "林書逸", "手忙腳亂"},
		},
		{
			FileName: "004.html",
			GUID:     "news-politics-paper-1600975",
			Category: "政治",
			Author:   []string{"吳正庭", "甘孟霖"},
			Title:    "侯、郭金門較勁 地方政要選邊站",
			Tags:     []string{"侯友宜", "郭台銘", "823砲戰", "陳玉珍", "2024總統大選"},
		},
	}

	for i := range tcs {
		tc := tcs[i]

		t.Run(
			tc.FileName,
			func(t *testing.T) {
				// t.Parallel()
				f, err := os.Open(fmt.Sprintf("example/ltn/%s", tc.FileName))
				require.NoError(t, err)
				q := p.Parse(parser.NewTestQuery(200, "", f))
				require.NoError(t, q.Error)
				require.NotNil(t, q)
				require.Equal(t, tc.GUID, q.News.GUID)
				require.Equal(t, tc.Title, q.News.Title)
				require.Equal(t, tc.Category, q.News.Category)
				require.False(t, q.News.PubDate.IsZero())
				require.ElementsMatch(t, tc.Author, q.News.Author)
				require.ElementsMatch(t, tc.Tags, q.News.Tag)
				if len(tc.Content) > 0 {
					require.Equal(t, len(tc.Content), len(q.News.Content))
					for i, c := range q.News.Content {
						require.Equal(t, tc.Content[i], c)
					}
				}
			},
		)
	}
}

func TestUP(t *testing.T) {
	p := parser.NewUPParser()

	type testCase struct {
		FileName string
		Title    string
		Category string
		Author   []string
		GUID     string
		Tags     []string
		Content  []string
	}

	tcs := []testCase{
		{
			FileName: "001.html",
			Title:    "海葵「路徑南修」恐直擊北台灣　預估今轉中颱並發布海警、明陸警",
			Category: "焦點",
			Author:   []string{"上報快訊"},
			GUID:     "24-181114",
			Tags:     []string{"海葵", "颱風", "中央氣象局", "中颱", "登陸", "陸警", "海警"},
			Content: []string{
				"輕颱海葵原本預估從台灣北部海域掠過，但其路徑不斷南修，最新預測恐從台灣東北部登陸，可能成為睽違4年再度登陸台灣的颱風。中央氣象局觀測，海葵正緩慢增強中，可能在今天（1日）轉為中颱，最快可能下午就會發布海上颱風警報、明天清晨發布陸上颱風警報。",
				"中央氣象局監測，今年第11號颱風海葵，目前位置在台灣東方海面上空，時速約19公里，持續朝西北西移動。7級暴風半徑達120公里。原本海葵的路徑被預測比較偏北，但各國預測持續南修，今天最新預測，海葵有很大的機率登陸台灣本島。",
				"中央氣象局預估，最快在今天下午就會針對海葵發布海上颱風警報，今天深夜到明天（2日）清晨可能會發布陸上颱風警報，至於登陸的時間、地點仍待觀察。而海葵的強度也在增強中，應該會在今天下午成為中颱。",
				"國家災防科技中心（NCDR）則預警，海葵最接近台灣的時間為明天至後天（3日），台灣北部、東北部將會是受影響最劇烈的地區，而中部以北則是都會降雨，提醒民眾需做好防颱準備。",
			},
		},
		{
			FileName: "002.html",
			Title:    "名廚江振誠餐廳「RAW」沒摘米其林三星　IG全黑圖原因曝光",
			Category: "焦點",
			Author:   []string{"鍾知諭"},
			GUID:     "24-181100",
			Tags:     []string{"江振誠", "米其林", "三星", "名廚", "RAW"},
			Content: []string{
				"國際名廚江振誠開設的餐廳「RAW」，連續多年拿下台北米其林二星，被網友認為是「全台最難訂餐廳」，不過「台灣米其林指南2023」星級餐廳昨天（8月31日）頒獎，RAW卻沒有獲得三星，江振誠隨即在IG發出全黑圖，引發臆測，他晚間在臉書說明，「因為表妹離開10年，很想她。」",
				"「台灣米其林指南2023」星級餐廳昨天公布入選三星的餐廳包括6度蟬聯三星的「頤宮中餐廳」，還有都是從二星升級為三星的「JL STUDIO」和「態芮TaÏrroir」。",
				"至於江振誠的餐廳RAW沒有晉升三星，一樣維持在二星榜單。江振誠昨天在IG上發出全黑圖，且沒有留下任何話。全黑圖的貼文也引起網友關注，認為他是為了米其林不開心，還有不少人替他打抱不平。",
				"不過，江振誠隨即發出2篇貼文，表示自己在新加坡，還要繼續亞洲巡演，又提到自己是因為表妹離開10年了，很想她、沒什麼好分享的，大家成熟一點。",
			},
		},
		{
			FileName: "003.html",
			Title:    "【有片】華為低調推新機　中國官媒卻嗨喊「打贏對美科技戰」",
			Category: "國際",
			Author:   []string{"洪毅"},
			GUID:     "3-181068",
			Tags:     []string{"華為", "美中科技戰", "中國", "美國", "手機", "晶片"},
			Content: []string{
				"中國手機大廠華為29日無預警推出最新款手機Mate 60 Pro，雖然仍有許多細節不明朗，但已激起部分中國民眾的民族情結，稱這是趁著美國商務部長雷蒙多訪問中國期間的有意為之，中國官媒「環球時報」也發布社論，要雷蒙多聽到在美國打壓之下「昂起頭來的聲音」。",
				"彭博報導，華為沉寂3年後，29日推出新品Mate 60 Pro讓中國官媒掀起一股民族自豪感，將這款新機形容為科技奇蹟，在美國制裁中贏得勝利。環球時報直指，作為被美國「卡脖子」最厲害的中國科技公司 ，華為經歷艱難的低谷之後，能否解決沒有可搭載5G技術晶片的問題，是中國科技業承受壓力的重要標誌。不少民眾認為，華為在雷蒙多（Gina Raimondo）訪中期間推出新機並非巧合，因為美國商務部是執行制裁的機構，至今仍有數百家中國企業在其「實體清單」上。",
				"環球時報在標題為「雷蒙多該如何準確理解華為預售新機」的社論中指出，中國網壇對華為新機的興奮之情，「凸顯出對我國科技自主研發的强烈期待和信心」，許多網友認為，華為發表新機的時間點有「在美國打壓之下昂起頭来」的深層次異議，這些聲音應當被雷蒙多以及更多美國人聽到，並對華盛頓形成觸動。",
				"社論還說稱，美國對中國發動貿易戰「毫無疑問已經被證明失敗了」，科技戰仍在進行中，但是華為推出Mate 60 Pro，證明美國打壓已經失敗，預示科技戰的最終結果，中國「在打壓之下奮起直追的勁頭，以及支撑這股信念的強大道義感，卻是美國比不過的」，希望雷蒙多能把華為新機引發的熱情帶回美國，「並在一定程度上讓美國更加能夠真正讀懂中國」。",
				"網傳華為的新機具有5G功能，對於Mate 60 Pro處理器的來源與性能，華為拒絕置評，只表示不會在中國以外的地方推出。新浪新聞報導，華為的Mate 60 Pro是直接在網路上推出，開賣算是較為低調，雖然外界十分關注此手機究竟使用什麼晶片，華為僅稱晶片能強化手機的通信體驗及網路連接，未透露太多細節。",
			},
		},
		{
			FileName: "004.html",
			Title:    "《知否》朱一龍博命演出災難片《峰爆》　冒失溫危險拍出「這1幕」大獲好評",
			Category: "流行",
			Author:   []string{"李雨勳"},
			GUID:     "196-180968",
			Tags:     []string{"知否？知否？應是綠肥紅瘦", "朱一龍", "人生大事", "峰爆", "消失的她", "陸劇", "陸片", "金雞獎", "Netflix"},
		},
		{
			FileName: "005.html",
			Title:    "【蘇拉襲台】小三通「金門往返五通」明起全天停航　交通、活動異動一次看",
			Category: "焦點",
			Author:   []string{"鍾知諭", "袁維駿"},
			GUID:     "24-180829",
			Tags:     []string{"颱風", "蘇拉", "航空", "船運", "交通部"},
		},
		{
			FileName: "006.html",
			Title:    "陳嘉宏專欄：柯文哲的心靈雞湯",
			Category: "評論",
			Author:   []string{"陳嘉宏"},
			GUID:     "2-181092",
			Tags:     []string{"柯文哲", "心靈雞湯"},
			Content: []string{
				"柯文哲說，要知道現在的選情如何，看民進黨鎖定誰在進行攻擊就知道了。這句話不能說有錯，但其實不夠精準。要判讀選情如何？聽柯文哲怎麼講話其實是最準的。如果你常常看到柯文哲張牙舞爪、大放厥詞，那就是柯文哲的民調高的時候，因為他得意到屁股都翹起來了；如果你看到柯文哲欲言又止，正經八百地想跟你談政策，那就是他開始想轉移話題，改變風向，表示他民調不妙了！",
				"例如，早前柯文哲罵人是「狗」，指控北市府官員是「太監」，辦了一場不知所以的「募款演唱會」，找來假空姐舞團，都正是他民調高檔的時候。但這兩個禮拜以來，柯文哲突然態度丕變，宣稱要找政黨領袖對話，還弄了一個什麼「四葉草運動」說要推動「世代和解」，想當然是他民調往下走了。用柯文哲的嘴巴來觀察政治方向非常準，比天上的北極星還準。",
				"柯文哲藏不住話，外界常常可以從他講的話知道他的思考模式。例如柯文哲說，他拿過杜聰明的獎學金，「當時醫學系第一名，我拿到那個獎很大，大概可以一個學期的註冊費。」沒想到真的有人去翻出當時杜聰明基金會給的獎學金，不但獎學金得主另有其人，醫學系第一名也不是他。原來，柯文哲口中的獎學金不是基金會所頒發的正式獎學金，而是杜家給他們兒孫輩的私下獎學金（柯文哲算是杜聰明的姻親晚輩）。柯文哲可能沒說謊，但前述一番話實在也有魚目混珠之嫌。",
				"喜歡魚目混珠正是柯文哲的講話風格。例如他上次開「演唱會」時吹牛說負責伴奏的鍵盤手是五月天的，立刻惹來本尊出面否認，但原來柯文哲說的是以前的鍵盤手不是現在的鍵盤手。例如他說日本自民黨的幹事長跟他說台灣不可能加入CPTPP，原來柯文哲說的不是號稱「黨三役」的自民黨幹事長茂木敏充，而是參議院幹事長世耕弘成，一番話幾乎引來台日外交的軒然大波，他卻從頭到尾裝沒事。",
				"柯文哲不知道自民黨幹事長與參議院幹事長不一樣嗎？他不知道五月天過去的鍵盤手與現在的鍵盤手是不同一個人嗎？他不曉得杜聰明基金會的正式獎學金與杜聰明給兒孫輩的私下獎學金有不同意義嗎？柯文哲智商157，他怎麼可能不知道這其中的差別。知道還這樣混淆視聽，其實就是因為他生性喜歡佔點小便宜；他特別享受別人對他的「崇拜」，所以往往「做嘸一湯匙」卻「講得一畚箕」。即使已經快65歲領老人卡了，他的行徑還是像個中二生，愛談當年勇，到處跟人家比成績比學校比科系。",
				"外界也很容易從柯文哲的話語裡讀出他的意向思維，例如，他在2014年是反核急先鋒，還曾說：「核四的問題沒那麼困難，看你可考慮的是未來4年、未來10年，還是未來100年。」但最近他要選總統卻說道：「如果你要台積電，那你就要想想看你要不要核電。」明眼人都知道，「若要台積電」是假，「沒辦法不接受核四」才是真的。只是柯文哲當年反核的印記太深，所以用「接受台積電」來當他的遮羞布。",
				"台積電的確是台灣的用電大戶，不過，台積電用的電其實與核電毫無干係，因為台積電早在2020年就加入了RE100，而RE100只接受綠電，並不接受核電。柯文哲拿一項錯誤資訊到處放送，無非是急著想透過支持核電進取藍營選票。在選舉思維下，反核不再重要，他口中「未來100年」的事也不再重要，只有他的選舉最重要。只要耐心地對照柯文哲的前言後語，當知道他是一個沒有中心思想，隨時可以出賣自己價值的政客。",
				"柯文哲當過台大外科急診中心的主任，每次演講總是以「已經看透生死」破題。他說：「我要建立公平社會，讓年輕人只要努力就有機會。」「理性科學務實來處理政治問題」「面對兩岸關係就是，可以合作的時候合作，應該競爭的時候競爭，必須對抗的時候就要對抗。」這些話乍聽之下莫測高深，但其實只是一碗碗心靈雞湯。心靈雞湯可用來調節個人心性，卻絕無任何政策的可行性；若像過去幾年的台北市民這樣不辨滋味、囫圇吞棗，勢必害人也害己。",
			},
		},
	}

	for i := range tcs {
		tc := tcs[i]

		t.Run(
			tc.FileName,
			func(t *testing.T) {
				t.Parallel()
				f, err := os.Open(fmt.Sprintf("example/up/%s", tc.FileName))
				require.NoError(t, err)
				q := p.Parse(parser.NewTestQuery(200, "", f))
				require.NoError(t, q.Error)
				require.NotNil(t, q)
				require.Equal(t, tc.GUID, q.News.GUID)
				require.Equal(t, tc.Title, q.News.Title)
				require.Equal(t, tc.Category, q.News.Category)
				require.False(t, q.News.PubDate.IsZero())
				require.ElementsMatch(t, tc.Author, q.News.Author)
				require.ElementsMatch(t, tc.Tags, q.News.Tag)

				if len(tc.Content) > 0 {
					require.Equal(t, len(tc.Content), len(q.News.Content))
					for i, c := range q.News.Content {
						require.Equal(t, tc.Content[i], c)
					}
				}
			},
		)
	}
}

func TestRFI(t *testing.T) {
	p := parser.NewRFIParser()

	type testCase struct {
		FileName string
		Title    string
		Category string
		Author   []string
		GUID     string
		Tags     []string
		Content  []string
	}

	tcs := []testCase{
		{
			FileName: "001.html",
			Title:    "超強颱風「蘇拉」襲粵 深港澳停工停課 今晚風更勁 逾400航班受影響",
			Category: "港澳台",
			Author:   []string{"香港特約記者 麥燕庭"},
			GUID:     "qbc+7tL4Ge3fQx+j7BTy",
			Tags:     []string{"港澳台", "中國", "香港", "自然災害"},
			Content: []string{
				// "被形容為「完美風暴」的超強颱風「蘇拉」，連日令兩岸四地政府嚴陣以待，在沒有造成傷亡地略過台灣後，便逼近香港，令港府破例地提早半天預告，今(1日)天凌晨會發出第三高的 8號熱帶氣旋警告信號，當颱風傍晚最接近香港時更會考慮是否改掛更高風球。截至中午，已有部分地區出現水浸，居民要撤離；海陸空交通亦大受影響，暫有三百多班航班要取消，下午二時後，所有航班將基本上停航。",
				"深圳「五停」 廣東最高級別戒備",
				"中國中央氣象台亦因應「蘇拉」發布最高級別的颱風紅色預警，預測該颱風將於今天下午至夜間在廣東省惠來至香港一帶沿海登陸。面對蘇拉可能是1949年建國以來登陸廣東的其中一股最強颱風，廣東省已將防颱風應急響應提升至最高一級；深圳市亦發出緊急動員令，宣布全市下午逐步「五停」：停課丶停工丶停業丶停市丶停運，而深圳機場更是在今天中午起已暫停航班升降，當局更要求居民，如非必要，不要外出；至於廣東汕尾更是今早已全市停工丶停課丶停運，並疏散近一萬人到安全地方，近八千艘各類漁船已回港避風。",
				"作為防備，中海油位於南海的油氣田，疏散七千多名海上作業人員。另外，中國國家防汛抗旱總指揮部已加派工作組到廣東福建協助，交通運輸部救助打撈局亦已部署16艘大功率救助船和 9架救助直升機加強應對。",
				"港府總動員 下午航班停 今晚考慮改掛 8號以上風球",
				"香港方面，政府亦嚴陣以待，以「全政府動員」規格迎接超強颱風「蘇拉」，意味有需要時可動員占港府編製約 55%的至少一萬人處理及善後。作為港府第二把手的政務司司長陳國基昨午便率領10個局長級和處長級官員舉行記者會，公布應對「蘇拉」可能為香港帶來的嚴重威脅的應變部署。",
				"天文台長陳栢緯表示，「蘇拉」在香港部分地區引發的風暴潮有可能達至歷史高位，沿岸低窪地區或嚴重水浸，若「蘇拉」在香港以南掠過，情況或會與 2018年「山竹」襲港相若，預計吐露港丶維港及大澳將會成為重災區，其中吐露港水位或達海圖基準面上5米，接近1962年「溫黛」襲港紀錄。",
				"翻查紀錄，「溫黛」襲港期間，天文台懸掛最高的 10號風球，風暴期間有 183人死亡或失蹤，是二次世界大戰後造成最多人員傷亡的颱風。更令人提高警戒的是，「溫黛」當年亦是在 9月1日襲港，天文台總部當時錄得最高 60分鐘平均風速達到每小時 133公里；最高陣風每小時 259公里；最低瞬時海平面氣壓 953.2百帕，三項紀錄至今仍未被打破。",
				"至於 2018年9月16日襲港的超強颱風「山竹」， 在香港造成超過 360人受傷，多個公共交通網絡陷於癱瘓，港府因為沒有宣布停工而備受抨擊。",
				"天文台自凌晨接近3時已因「蘇拉」發出 8 號熱帶氣旋警告信號，預計這個俗稱的 8號風球會在今天餘下大部分時間懸掛，天文台中午時預計，按照現時預測路徑，「蘇拉」會在下午5時以每小時 185公里的中心最高風速進入香港100公里範圍內，並於今晚至明早最接近香港，在天文台以南約 50公里內掠過。",
				"天文台續稱，會視乎本地風力變化，考慮在晚上 6時至 10時之間發出更高熱帶氣旋警告信號。至於是改掛 9號抑或最高的10號風球，天文台未有透露。若預測不變，蘇拉到周日基本上會遠離香港，周日及周一風勢會減弱，但仍會有驟雨。",
				"當 8號風球高懸時，全港已停工停課，海上和陸上交通亦續步停駛，航班則在可行情況下維持，機場管理局表示，截至中午，有大約 366班航班取消丶約40班航班延誤，而順利運作的航班有 600班。至於下午 2時後的航班，基本上會取消，但機場不會關閉。",
				"由於有超過三百班航班取消，大批旅客到場了解情況，以致航空公司櫃位大排長龍。另外，亦有乘客因為資訊有誤，錯過昨晚的班機，預計須在機場乾等一天，改乘今晚的航班離開，現時不知航班會否停航，需要再到公司櫃位了解。",
				"澳門停通關 晚上改掛9號風球機會高",
				"澳門方面，當地氣象局亦於今日下午 2時改發 8號西北風球，預測當地風力將會增強，而今晚至明日凌晨改發 9號風球機會是中等至較高，清晨改發10號風球機會中等。氣象局續稱，距離「蘇拉」中心外圍約100至120公里附近，已錄得10級或以上風力，預料「蘇拉」有機會正面吹襲澳門。",
				"氣象局又指，由於蘇拉將會非常接近珠江口，受天文大潮疊加風暴潮影響，預料明早低窪地區有機會出現 1.5米左右水浸，水位可能在短時間內快速上升。",
				"受「蘇拉」影響，連接澳門和珠海橫琴的蓮花大橋將會封閉，而橫琴口岸亦會在下午3時許暫停通關。",
			},
		},
		{
			FileName: "002.html",
			Title:    "澤連斯基稱烏克蘭製造的武器已被證明射程達700公里",
			Category: "國際",
			Author:   []string{"弗林"},
			GUID:     "CVs3jX4lUxv6ruezvPIS",
			Tags:     []string{"國際", "烏克蘭", "俄羅斯"},
			Content: []string{
				"據當地報道，莫斯科周三指責烏克蘭的無人機對俄羅斯與愛沙尼亞和拉脫維亞邊境附近的機場進行了長達4小時的襲擊，並造成4架伊爾-76軍用運輸機受損。該機場位於俄羅斯普斯科夫州，距烏克蘭邊境以北約700公里處。在長達18個月的戰爭中，俄羅斯共有6個州成為襲擊目標。",
				"澤連斯基周四晚在發表例行講話中談到，“今天是多事的一天。與軍方和政府官員舉行電話會議。前線。我們的進攻行為。我們的武器，即烏克蘭的新型武器，其射程為700公里。任務就多了”。他當天還與英國貝宜系統公司的代表進行了會談。",
				"澤連斯基介紹說，“全世界都非常了解這家公司。我們的戰士已經非常熟悉該公司生產的武器。特別是火炮——L119和M777榴彈炮及CV90步兵戰車，威力非常大。該公司開始在烏克蘭開展業務。我們的目標是所有最有用的防禦武器都可以在烏克蘭生產。我們已經生產了一些產品，我們將生產所有必要的產品。我感謝世界上所有提供幫助的人！”",
				"澤連斯基近日在接受葡萄牙廣播電視台採訪時表示，除了盟國已經承諾的“50到60架”美製F-16戰鬥機外，烏克蘭還需要另外100架戰鬥機，這些戰鬥機應在明年年初投入使用。",
				"澤連斯基在採訪中說，“今天，我們就未來交付50至60架（F-16）戰鬥機達成了協議。我為什麼這麼說？因為在不同時期，我們將擁有不同數量的戰鬥機”。他補充說：“我們總共需要約 160架戰鬥機，才能擁有一支強大的空軍，不給俄羅斯主宰領空的機會”。",
				"澤連斯基表示，這些戰鬥機應該“明年初”在烏克蘭天空投入使用，他認識到這個問題很複雜，因為不僅需要對烏克蘭飛行員進行培訓，還需要對該國工程師進行培訓，並進行大量的專業維護。",
				"澤連斯基表示，“我們正在與俄羅斯作戰，我們正在為我們的烏克蘭土地而戰，反對俄聯邦的侵略政策。我們需要戰鬥機只是為了保衛我們自己。保衛我們的土地、我們的海洋、我們的天空”。他指出，烏克蘭需要西方戰鬥機的另一個緊迫用途是確保“俄羅斯不會非法統治黑海，會通過封鎖我們在黑海和亞速海的糧食走廊來控制俄羅斯的侵略”。",
			},
		},
	}

	for i := range tcs {
		tc := tcs[i]

		t.Run(
			tc.FileName,
			func(t *testing.T) {
				t.Parallel()
				f, err := os.Open(fmt.Sprintf("example/rfi/%s", tc.FileName))
				require.NoError(t, err)
				q := p.Parse(parser.NewTestQuery(200, "", f))
				require.NoError(t, q.Error)
				require.NotNil(t, q)
				require.Equal(t, tc.GUID, q.News.GUID)
				require.Equal(t, tc.Title, q.News.Title)
				require.Equal(t, tc.Category, q.News.Category)
				require.False(t, q.News.PubDate.IsZero())
				require.ElementsMatch(t, tc.Author, q.News.Author)
				require.ElementsMatch(t, tc.Tags, q.News.Tag)

				if len(tc.Content) > 0 {
					require.Equal(t, len(tc.Content), len(q.News.Content))
					for i, c := range q.News.Content {
						require.Equal(t, tc.Content[i], c)
					}
				}
			},
		)
	}
}

func TestTVBS(t *testing.T) {
	p := parser.NewTVBSParser()

	type testCase struct {
		FileName string
		Title    string
		Category string
		Author   []string
		GUID     string
		Tags     []string
		Content  []string
	}

	tcs := []testCase{
		{
			FileName: "001.html",
			Title:    "海葵會變強烈颱風「紅色警戒區壟罩全台」　日本預測路徑曝",
			Category: "全球",
			Author:   []string{"陳妍如"},
			GUID:     "world-2227667",
			Tags:     []string{"颱風海葵", "台灣", "登陸", "海葵颱風", "颱風專區", "日本氣象廳", "沖繩", "海葵颱風路徑"},
			Content: []string{
				"今年第11號颱風「海葵」將，根據日本氣象廳公布的最新預測路徑，幾乎全台灣都被海葵颱風的紅色暴風圈警戒區壟罩；日本氣象協會也指出，目前海葵颱風位於沖繩南方海域，週六會靠近宮古島南方，而預計周日（3日）上午就會登陸台灣。",
				"日本《NHK》報導，海葵颱風正在沖繩以南海域朝西前進，今後將持續增強並逐漸靠近日本沖繩縣的先島群島（宮古群島、八重山群島）。先島群島從今（1日）起、沖繩本島從明（2日）起將出現大浪。",
				"專家指出，海葵颱風會繼續增強，預計會達到日本氣象廳所定義的「強烈颱風」，最大風速介於每秒33公尺至44公尺。根據預測，台灣本島幾乎全部處於每秒風速25公尺以上的紅色暴風區警戒範圍。",
				"此外，由於目前日本本州附近的太平洋高壓與海葵颱風之間的氣壓差仍在不斷增大，因此海葵颱風北側的強風區預計也會擴大，將會為沖繩本島、鹿兒島縣的奄美地區帶來強風、強降雨。",
				"據日本氣象廳1日上午9時（台灣時間8時）公布的數據，海葵颱風目前位於沖繩群島南方海域，中心位置北緯22度、東經129.4度，正以每小時15公里的速度朝西北西方向前進，中心氣壓980百帕，近中心最大風速每秒35公尺，最大瞬間風速每秒50公尺。",
				"海葵颱風直撲北台而來，《TVBS新聞網》9/1-9/3連續3天，每天16:00~22:00，都有颱風特別報導，跟主播互動分享各地風雨狀況。陪您掌握海葵颱風最新動態，敬請鎖定。",
				"直播連結：https://parser.tvbs.com.tw/live/news4live/34664",
			},
		},
		{
			FileName: "002.html",
			Title:    "有望「連放3天假」？海葵進補變更強　日預測：恐再南修撲雙北",
			Category: "生活",
			Author:   []string{"樓冠陞"},
			GUID:     "life-2227514",
			Tags:     []string{"颱風專區", "颱風海葵", "海葵颱風", "海葵颱風路徑", "海葵南修", "台灣颱風論壇｜天氣特急", "登陸", "日本氣象廳"},
			Content: []string{
				"今年第11號颱風「海葵」的路徑在一夕間「大幅西修」，不但可能成為「西北颱」外，颱風中心甚至不排除會從台灣東北角登陸。此外，氣象粉專「台灣颱風論壇｜天氣特急」也透露，海葵颱風路徑西修後，正好進入了「暖水池」區域，也意味著它極有可能在短時間內迅速增強。",
				"「台灣颱風論壇｜天氣特急」今（31）日稍早在臉書粉專發文指出，從太平洋最新的海水熱焓量分布圖可以發現，台灣以東的海面溫度目前高達30度，也意味著此區域的海水熱焓量較高、暖水層比較深厚，故可以提供颱風更充足的水氣與能量。",
				"該粉專解釋，海葵颱風接下來將會進入這塊「暖水池」區域，故在沒有其他因素的干擾之下，它將從中吸取豐沛的能量，並有機會在短時間內迅速成長。不過，好消息是海葵颱風的移動速度較快、吸收能量的時間自然也比較短，故巔峰強度預計只會落在「中度颱風中段班」。",
				"值得一提的是，根據日本氣象廳的最新預測，海葵颱風的未來路徑「再度南修」，不但可能登陸外，颱風中心更有機會於9月3日晚間9時，正好落在台北市上空。另外，日本氣象廳也預估，海葵颱風的移動速度將從目前的每小時20公里，漸漸放慢至每小時15公里，不排除會拉長影響台灣天氣的時間。",
			},
		},
		{
			FileName: "003.html",
			Title:    "3歲童突高燒！詭異水泡連成排　凶手竟是「這款鞋」",
			Category: "健康",
			Author:   []string{"樓冠陞"},
			GUID:     "health-2228412",
			Tags:     []string{"孩童", "蜂窩性組織炎", "膠鞋", "水泡", "發燒", "束帶", "鰻魚家家酒"},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.FileName,
			func(t *testing.T) {
				// t.Parallel()
				f, err := os.Open(fmt.Sprintf("example/tvbs/%s", tc.FileName))
				require.NoError(t, err)
				q := p.Parse(parser.NewTestQuery(200, "", f))
				require.NoError(t, q.Error)
				require.NotNil(t, q)
				require.Equal(t, tc.GUID, q.News.GUID)
				require.Equal(t, tc.Title, q.News.Title)
				require.Equal(t, tc.Category, q.News.Category)
				require.False(t, q.News.PubDate.IsZero())
				require.ElementsMatch(t, tc.Author, q.News.Author)
				require.ElementsMatch(t, tc.Tags, q.News.Tag)

				if len(tc.Content) > 0 {
					content := strings.Join(q.News.Content, "")
					for _, c := range q.News.Content {
						require.Contains(t, content, c)
					}
				}
			},
		)
	}
}

func TestSETN(t *testing.T) {
	p := parser.NewSETNParser()

	type testCase struct {
		FileName string
		Title    string
		Category string
		Author   []string
		GUID     string
		Tags     []string
		Content  []string
	}

	tcs := []testCase{
		{
			FileName: "001.html",
			Title:    "海葵遠離！台灣本島今晚20:30解除陸警",
			Category: "生活",
			Author:   []string{"賴俊佑"},
			GUID:     "1348570",
			Tags:     []string{"陸警", "海葵", "颱風"},
			Content: []string{
				"海葵颱風強度持續減弱，不過暴風圈仍籠罩雲林以南至高雄陸地及澎湖，氣象局表示颱風在今晚會持續朝西移動，預計台灣本島今晚解除陸上警報，不過宜花地區、東北角和金門受到颱風環流影響，仍有較大雨勢出現。",
				"氣象局資深預報員謝佩芸說明，海葵颱風今天下午在澎湖南方海面呈現滯留打轉的現象，今日下午5點颱風中心仍在澎湖西南方大約50公里海面上，向西北西轉西南西移動，在過去3小時強度略為減弱，暴風半徑維持120公里，其暴風圈仍籠罩雲林以南至高雄陸地及澎湖，此颱風暴風圈正逐漸接近金門地區，預計未來此颱風強度有持續減弱的趨勢，並持續針對台灣海峽以及東沙島海面發布海上警報",
				"氣象局表示，颱風在今晚會持續朝西移動，台灣本島預計今晚解除陸上警報，颱風中心預計明日清晨到上午會登陸中國福建沿海，預計明天上午台灣會解除陸上警報；颱風環流仍持續影響台灣，氣象局表示，預估今晚20點後雨勢趨緩，但宜花地區、東北角和金門仍有比較大的雨勢出現，",
				"另外，今天生成的熱帶低氣壓有機會在明天上午增強為輕度颱風，未來將會朝東北方移動，這個輕度颱風延伸到海葵颱風的位置，將會是一個低壓帶，台灣正位於這個低壓帶中間，天氣相當不穩定，",
				"氣象局預估，週二(9/2)持續受到海葵颱風外圍環流影響，各地斷斷續續有短暫陣雨，中午過後中部以北要留意局部大雨，週三至週五(9/3~9/5)台灣位於低壓帶間，全台須留意午後大雷雨，雨勢要到週六才會趨緩。",
				"海葵颱風對台灣的影響逐漸減弱，但提醒澎湖、金門的民眾要留意颱風帶來的風雨，未來幾天台灣天氣相當不穩定，提醒民眾外出要攜帶雨具並留意最新氣象訊息。",
			},
		},
		{
			FileName: "002.html",
			Title:    "寶可夢中心坐落信義區？偵探系網友翻證據",
			Category: "生活",
			Author:   []string{"林柏廷"},
			GUID:     "1348583",
			Tags:     []string{"信義區", "寶可夢中心", "A11", "皮卡丘"},
			Content: []string{
				`寶可夢台灣於1日在臉書粉專宣布「Pokémon Center TAIPEI」將於12月開幕，讓許多訓練家相當期待，沒想到有訓練家發揮柯南的精神，藉由室內裝修變更查詢，發現寶可夢中心的設立地點，很有可能就在信義區的A11。`,
				`全台首間Pokémon Center將於12月在台北盛大開幕，預計推出開幕限定紀念商品、舉辦開幕活動等，不過詳細的確切位置並沒有公開，令許多訓練家十分好奇。`,
				"不過有偵探系的訓練家根據室內裝修變更查詢，意外發現寶可夢中心似乎就坐落在信義區的A11，取代原先的無印良品A11門市。",
				`根據104人力銀行的「Pokémon Center TAIPEI 寶可夢中心兼職店員」招募頁面所述，上班的地點在台北市信義區，更加印證了坐落在信義區A11的傳聞。`,
			},
		},
		{
			FileName: "003.html",
			Title:    "拜習會生變！「他」代替習近平出席G20",
			Category: "國際",
			Author:   []string{"許元馨"},
			GUID:     "1348480",
			Tags:     []string{"中國", "G20", "習近平", "拜登"},
			Content: []string{
				"中國國家主席習近平將不出席在印度舉行的20國集團（G20）領袖峰會，對此，美國總統拜登表達感到失望。今（4）日，中國外交部發言人毛寧宣布，應印度共和國政府邀請，中國國務院總理李強將於9月9日至10日出席在印度新德里舉行的二十國集團領導人第十八次峰會。",
				"根據路透社報導，針對習近平不出席G20峰會，拜登在德拉瓦州雷荷勃斯灘（Rehoboth Beach）表示：「我感到失望…但是我會碰得上他。」但拜登未就此詳述。",
				"中國外交部發言人毛寧今日宣布，應印度共和國政府邀請，中國國務院總理李強將於9月9日至10日出席在印度新德里舉行的二十國集團領導人第十八次峰會。",
			},
		},
		{
			FileName: "004.html",
			Title:    "憂吃輻射海鮮買檢測儀！陸女測大閘蟹嚇傻",
			Category: "兩岸",
			Author:   []string{"CTWANT"},
			GUID:     "1347793",
			Tags:     []string{"中國", "海鮮", "輻射", "大閘蟹", "檢測儀", "核處理水", "日本"},
			Content: []string{
				"日本核處理水排海事件引發輿論關注，在各大電商平台，核輻射檢測儀成為暢銷爆款。不過，大陸上海市一名徐姓女子反映，8月29日晚間她在家中蒸了一鍋大閘蟹，拿出一個月前在電商平台上購買的核輻射檢測儀，結果警報聲不斷，頻繁警示劑量超標，令她和家人感到十分恐慌。",
				"《上觀新聞》報導，徐女反映，她日前在家中蒸了一鍋大閘蟹，拿出一個月前在電商平台上購買的核輻射檢測儀，結果警報聲不斷，頻繁警示劑量超標。更諷刺的是，徐女隔天又拿檢測儀對準女兒肚皮，警報聲再次響起，讓她忍不住懷疑這款號稱精準的核輻射檢測儀也許名不符實，於是與線上客服進行交涉。",
				"客服表示，可能是個別產品的質量問題，可以補發貨或者退款，但拒絕電話溝通。實測該款檢測儀，在桌面上靜置約2分鐘，顯示數據飆升至每小時3.33微西弗，達報警值3倍以上。徐女撥打客服熱線，客服人員讓她將視頻發過去，並安排專人核實，因為該產品沒辦法直接走退換流程，等把問題核實清楚，再來協調方案。",
				"對此，專家指出，放射性元素測量的前置條件較為複雜，通常只用一個儀器無法完成，且檢測方式也有嚴格的要求，因此普通百姓購買放射性檢測儀必要性不大，且針對海洋產品去測試核輻射，以前從未發生過，實在難以證實廠家的產品是否有作用，再說民眾只是為了求安心，不需要把自己搞得更不安心。",
			},
		},
		{
			FileName: "005.html",
			Title:    "防颱奇招！餐廳老闆砸12萬　租砂石車圍店",
			Category: "兩岸",
			Author:   []string{"李育道"},
			GUID:     "1347825",
			Tags:     []string{"玻璃", "砂石車", "防颱", "餐廳"},
			Content: []string{
				"颱風海葵直撲台灣今傍晚於台東地區登陸可能性很高，有機會終結自2019年8月24日以來「長達1471天沒有颱風登陸的狀況」，也因此不少民眾已在提前做好防颱準備；而今年第9號颱風蘇拉已於昨（2）日的凌晨3時30分以強颱風級在廣東珠海市金灣區沿海登陸，汕尾市就有一家海景餐廳為了防颱，花費3萬多元人民幣（約12萬元台幣），找來12輛砂石車包圍餐廳，以免餐廳被強風豪雨破壞。中國《星視頻》報導，餐廳林姓老闆表示，自家餐廳去年才開業，距離海邊只有100多公尺，由於餐廳內大面積使用落地玻璃讓顧客能夠欣賞海景，因此如果被颱風破壞，恐怕損失慘重，與餐廳受損相比，寧願租砂石車擋風雨，來保護餐廳「開這2萬相對便宜」。",
				"林老闆透露，其實上次中颱杜蘇芮登陸時，就曾租砂石車保護餐廳，而這次更加謹慎，還加錢要求砂石車裝土增重，讓防禦力提升，最後成功避免餐廳受到損害。",
				"台灣部分，中颱海葵強襲，氣象局持續發布海上、陸上警報，根據今（3）日凌晨最新觀測，海葵5時暴風圈已觸陸，暴風圈正逐漸進入臺灣東部及東南部陸地，對宜蘭、花蓮、臺東、新竹以南地區、恆春半島及澎湖構成威脅，預估將在下午登陸台東縣，今日到明日是對台灣影響最劇烈的時候。提醒民眾最好防颱準備。",
			},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.FileName,
			func(t *testing.T) {
				t.Parallel()
				f, err := os.Open(fmt.Sprintf("example/setn/%s", tc.FileName))
				require.NoError(t, err)
				q := p.Parse(parser.NewTestQuery(200, "", f))
				require.NoError(t, q.Error)
				require.NotNil(t, q)
				require.Equal(t, tc.GUID, q.News.GUID)
				require.Equal(t, tc.Title, q.News.Title)
				require.Equal(t, tc.Category, q.News.Category)
				require.False(t, q.News.PubDate.IsZero())
				require.ElementsMatch(t, tc.Author, q.News.Author)
				require.ElementsMatch(t, tc.Tags, q.News.Tag)

				if len(tc.Content) > 0 {
					require.Equal(t, len(tc.Content), len(q.News.Content))
					for i, c := range q.News.Content {
						require.Equal(t, tc.Content[i], c)
					}
				}
			},
		)
	}
}

func TestEtToday(t *testing.T) {
	p := parser.NewEtTodayParser()

	type testCase struct {
		FileName string
		Title    string
		Category string
		Author   []string
		GUID     string
		Tags     []string
		Content  []string
	}

	tcs := []testCase{
		{
			FileName: "001.html",
			Title:    "快訊／侯友宜將有「重大宣布」！　她喊：大家拭目以待",
			Category: "政治",
			Author:   []string{"曾羿翔"},
			GUID:     "20230904-2575420",
			Tags:     []string{"總統大選", "侯友宜", "張麗善"},
			Content: []string{
				"雲林縣長張麗善4日深夜表示，明天（5日）侯友宜會有重大宣布，請大家拭目以待，「我對侯友宜市長越來越有信心，他說到做到信守承諾，朋友們！明天等著他再提出好政見囉！」",
				"張麗善在臉書發文指出，侯友宜這幾次來訪雲林，給他的建議他都牢牢記住，而且化作當上總統後，會一一履行的政見。證實了侯深入基層體恤民意，願意解決人民的痛苦。國家的政策不公、財政收支劃分不公造成貧富不均及福利不平等，為什麼六都可以做的，其他縣市無法做到，福利制度過度懸殊，人口版圖自然移動，都往六都移動去了，造成嚴重城鄉差距。",
				"張麗善說，我告訴侯友宜「雲林拒作二等公民」，希望在台灣這塊土地上能得到平等尊嚴的對待，「相信其他縣市跟我一樣的期待，他聽進去了，我好感動！我對侯友宜市長，越來越有信心，他說到做到且信守承諾，朋友們！明天等著他再提出好政見。」",
				"侯友宜本人則在貼文底下留言，聆聽民眾的聲音，解決民眾的需求，是執政者永遠要牢記在心的事，既然民進黨做不到、不想做，我們就讓無能的政府下台，讓侯友宜跟大家一起，帶來改變，帶給民眾更好的日子！",
				"網友也紛紛留言表示，「加油」、「大家要團結」、「唯一支持侯友宜、唯一支持侯總統」、「侯友宜加油」、「堅定支持侯友宜」、「好期待」。",
			},
		},
		{
			FileName: "002.html",
			Title:    "超狂！大一新生從西藏走到四川報到入學　歷時2個月徒步上千公里",
			Category: "大陸",
			Author:   []string{"CTWANT"},
			GUID:     "20230904-2575147",
			Tags:     []string{"周刊王"},
			Content: []string{
				"9月是各大專院校的開學季節，住在遠地的新生也陸續報到，不過大陸一名大一新生，竟然靠著雙腿徒步走過上千公里，在歷時2個月的時間，靠雙腿走到四川托普信息技術職業學院報到。消息曝光後，引起不少網友的讚嘆與討論。",
				"四川托普信息技術學院日前在官方抖音帳號上發布一段影片，指稱一位名叫「三木旦·旺杰」大一新生，竟然花費暑假2個月的時間，徒步從西藏走到四川來做新生報到。三木旦·旺杰受訪時表示，自己從小就對徒步、登山很感興趣，又剛好暑假有2個月的時間，所以決定從西藏徒步走到四川的學校作新生報到。",
				"三木旦·旺杰表示，自己在過程中也經歷了滑坡、大雪、大雨，過程中也搭上一些便車、摩托車，甚至連拖拉機都坐過。而在這段過程中，三木旦·旺杰表示自己不僅認識了很多人，也見識到了過去從未見過的風景。",
				"不少網友在看到這段影片後，紛紛留言表示「他能從西藏走到四川，可我連床上到廁所都要猶豫好久」、「算一算時間，剛收到錄取通知書就要出發了」、「足根不痛嘛？我感覺4萬步是普通人極限了」、「內心勇敢，堅強，執行力強，動手能能力強，身體素質強，適應能力強，還有什麼，我暫時想不到了」、「有時候感覺古人真的很厲害，在沒有網路的時候就開始浪跡天涯了」。",
			},
		},
		{
			FileName: "003.html",
			Title:    "日本國會議員「福島衝浪比讚」　小泉進次郎：核處理水很安全",
			Category: "國際",
			Author:   []string{"鄒鎮宇"},
			GUID:     "20230904-2575377",
			Tags:     []string{"國會議員", "福島", "衝浪", "小泉進次郎", "核處理水", "日本"},
			Content: []string{
				"日本在8月24日將福島核廢水排入海中，整個排放期預計30年，周邊國家密切關注。對此，日本前環境大臣、現任國會議員、前首相小泉純一郎次子小泉進次郎4日到福島的海岸衝浪、吃生魚片，親自行動幫福島居民加油打氣。",
				"據日媒報導，小泉進次郎4日前往福島南相馬海岸烏崎海灘，參加一場為了兒童舉辦的衝浪課，並親自下場體驗衝浪，經過多次嘗試，成功在衝浪板上站起後也對著媒體比讚示意。",
				"小泉進次郎表示，日本內部及國外對於核處理水的評論不公正、無情也無科學根據，他親自體驗後不認同那些看法，就算衝浪會被說成是政治表演也沒關係，因為他就是來宣導核處理水排放入海後很安全，未來也會找其他的國會議員一起到福島衝浪、吃生魚片。",
				"許多日本網友看完紛紛留言，「這是一個很好的呼籲，只有進次郎才能做到」、「接下來換誰衝浪、釣魚或潛水」、「看他們玩得開心、吃得開心，我就有一種安全感」、「向你致敬，抱歉總是取笑你，我會反思的」、「我覺得這個人是天才，無法想像其他人能有這樣的表演」。",
			},
		},
		{
			FileName: "004.html",
			Title:    "快訊／台東縣宣布「2村、1校明停班課」　14學校明停課不停班",
			Category: "生活",
			Author:   []string{"楊漢聲", "趙蔡州"},
			GUID:     "20230904-2575371",
			Tags:     []string{"颱風", "海葵"},
			Content: []string{
				"颱風海葵過境，挾帶驚人雨勢，為台灣各地造成不少災情。台東縣政府4日晚間宣布，新增海端鄉初來國小新武分校5日停止上班、停止上課，成功鎮信義國小5日照常上班、停止上課，加上稍早已宣布的13所停課學校中，唯一停班停課的只有初來國小新武分校。",
				"台東縣政府稍早已宣布，海端鄉「霧鹿村」及「利稻村」因道路中斷，考量民眾安全，5日停止上班、停止上課。",
				"另外，東河鄉都蘭國中、泰源國中、泰源國小、都蘭國小、東河國小、北源國小、成功鎮三民國小及和平分校、鹿野鄉永安國小、瑞源國小、關山鎮電光國小、海端鄉海端國中、廣原國小等13所學校災損嚴重，考量學生安全，5日照常上班、停止上課。",
			},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.FileName,
			func(t *testing.T) {
				// t.Parallel()
				f, err := os.Open(fmt.Sprintf("example/ettoday/%s", tc.FileName))
				require.NoError(t, err)
				q := p.Parse(parser.NewTestQuery(200, "", f))
				require.NoError(t, q.Error)
				require.NotNil(t, q)
				require.Equal(t, tc.GUID, q.News.GUID)
				require.Equal(t, tc.Title, q.News.Title)
				require.Equal(t, tc.Category, q.News.Category)
				require.False(t, q.News.PubDate.IsZero())
				require.ElementsMatch(t, tc.Author, q.News.Author)
				require.ElementsMatch(t, tc.Tags, q.News.Tag)

				if len(tc.Content) > 0 {
					require.Equal(t, len(tc.Content), len(q.News.Content))
					for i, c := range q.News.Content {
						require.Equal(t, tc.Content[i], c)
					}
				}
			},
		)
	}
}

func TestBBC(t *testing.T) {
	p := parser.NewBBCParser()

	type testCase struct {
		FileName string
		Title    string
		Category string
		Author   []string
		GUID     string
		Tags     []string
		Content  []string
	}

	tcs := []testCase{
		{
			FileName: "001.html",
			Title:    "華為：新款5G手機橫空出世 中國半導體突破美國封鎖？",
			Author:   []string{"呂嘉鴻"},
			Category: "chinese-news",
			GUID:     "chinese-news-66748432",
			Tags:     []string{"中國", "中美關係", "技術新知", "經濟", "習近平", "華為", "貿易"},
			Content: []string{
				"中國電子巨擘華為在8月底無預警「低調」發售新5G智慧型手機Mate 60 Pro ，在全球半導體高科技界引發高度關注。",
				"目前各國科技公私單位，都在華為沒有正式宣佈手機零件細節前「拆機解謎」。多數結果發現，華為新機除了使用了中國半導體龍頭中芯科技研發的麒麟9000S七奈米晶片（芯片）之外，許多測試結果也顯示Mate 60 Pro運作速度可以與最新款的iPhone 5G手機一樣快。該機被認為是迄今為止中國本土技術生產的最先進版本。",
				"許多中國網民及若干官媒認為，該手機從配載的高科技7奈米晶片，到成為首款能與高軌衛星通訊的手機，是華為能夠突圍美國一系列制裁的證明，更象徵著中國半導體「創新自主」的未來指日可待。各地實體門市大排長龍，手機銷售一空。一些分析指，過去因受制裁而備受打擊的華為「起死回生」，為當下不安的中國經濟，注射了一劑愛國主義的強心針。",
				"但是，專家告訴BBC中文，華為此次與中芯國際合作生產的七奈米晶片，距離其他世界前沿晶片設計及代工的對手仍有一段距離。",
				"更重要的是，有專家認為，華為此次在美國商務部長雷蒙多拜訪北京的時機點上推出新機，政治信號濃厚，預計美國將對中國半導體祭出更強烈的制裁，並擴大審查與中國半導體及有關產業鏈。",
				"新加坡國立大學商學院高級講師卡普裏（Alex Capri）接受BBC中文採訪時說，此次華為手機在沒有發佈會和廣告的情況下低調銷售，又不公布手機零件細節且限量發佈，充滿「神秘氛圍」，「讓我感覺這是一場精心策劃的宣傳活動，也是中國廣泛的論述戰的一部分。」",
				"華為公關部通過郵件回復BBC中文的查詢，稱目前為止沒有太多聲明，並附上華為Mate 60 Pro新機發佈時的新聞稿，強調該手機是為全球消費者提供更高端及便捷的通訊。",
				"卡普裏說，無論如何，華為此次推出新機表明了中國在半導體行業上正全面性地加倍努力，以確保自己的國產晶片在設計和製造能力上朝著自主創新的方向前進。但他也強調，現在將華為Mate 60 Pro視為白宮對中國半導體制裁和出口管制的失敗還為時過早。",
				"美國將擴大施壓？",
				"9月5日，美國白宮國家安全顧問蘇利文（Jake Sullivan；沙利文）在記者會上回應說，在獲得有關華為新手機技術的組成資料前，不會對特定晶片和問題發表評論。",
				"但是，各方已經開始關注，美國在雷蒙多訪問之行意外收到中國來的「大禮」之後，美國是否會加大制裁範疇。根據《華盛頓郵報》，華為新機的推出在華府已造成擔憂，批評者稱，拜登對中國半導體的鉗製成效如何值得重新檢討。",
				"9月8日，彭博社報導美國商務部表示，將對華為新機的晶片及裝配展開調查。中國外交部發言人毛寧在同日回應此事稱，美國將商業行為政治化，但其「制裁、遏制、打壓阻止不了中國發展，只會增強中國自立自強、科技創新」的腳步。",
				"卡普裏告訴BBC，他預期白宮將更仔細地審查中芯國際對較老一代的美國或他國機器設備的使用情況，以及其他關鍵的產品輸入。",
				"台灣半導體評論者許美華說，接下來美國可能會對與美國有關的全球半導體產業鏈審查的越來越細，徹底盤整這些含有美國技術的全球半導體公司與中國的關係為何。她以台灣為例說，近日當地業界傳出有些台灣廠商，將半導體設備「提前報銷」之後，再轉售給中國廠商獲取利潤，這可能會是美方未來的關注方向之一。",
				"位於台北的智庫，台灣經濟研究院產經資料庫總監劉佩真則告訴BBC說，她預見美國對於中國將會更加防備，不論是去風險或是脫鉤的戰略，將持續浮上台面。「而中國受到華為此次事件的激勵，更將全力衝刺半導體國產化的進程，對於突破美國科技的封鎖線也有所寄托與信心，並視華為麒麟9000S晶片為中國反攻的灘頭堡」。在中美持續在半導體戰場上激戰下，全球半導體產業的「不確定性」（uncertainty），將持續存在。",
				"事實上，美中在半導體產業鏈上還有許多交集。美國科技諮詢公司（Moor Insights & Strategy） 高級分析師賽格（Ahsel Sag）告訴BBC說，譬如華為一直是美國科技大廠高通（Qualcomm）的大客戶，後者也曾經銷售5G技術給華為，直到今年初美國禁賣令發出後停止。華為與高通之間的技術關係，未來可能備受關注。",
				"無論如何，華為新機引發的國際關注，似乎證實了一個已經分裂和高度區域化的全球半導體貿易格局到來。卡普裏告訴記者，在半導體領域的脫鉤時代已經來臨了，「並將更廣泛地影響美中之間複雜的經貿關係。」",
				"華為Mate 60 Pro 橫空出世",
				"賽格解釋，此次華為新機最大的亮點該是5G數據機（modem；調製解調器)。這是考慮到華為已經近3年沒有在其手機中配備5G數據機了。",
				"此外，劉佩真分析，根據各單位拆解Mate 60 Pro手機及麒麟9000S晶片，此次華為供應鏈該是以中芯國際以DUV機台利用多重曝光達到七奈米制程來為其代工。電源方面則有聖邦、南芯半導體，無線充電則是美芯晟與封測技術的長電科技：「整體而言，華為主要是透過迂迴的採購與生產。」",
				"台灣南台科技大學朱岳中助理教授則告訴BBC，此次華為新機的最大賣點其實是衛星通訊。他表示，Mate 60 Pro是全球首款支援衛星通話的智慧型手機，這帶出兩個意義。首先是華為擺脫4G或5G的框架，採用衛星通訊，實測通訊速度更勝5G。他說，衛星通訊適用的電話並不稀罕，也不算貴，但體積通常較一般手機大，華為此款新機在於可以用標凖手機尺寸加入衛星電話功能。",
				"再者，有關衛星網路（如SpaceX公司的星鏈網路）目前為止都還要有額外接收器才能使用，但華為在手機中直接內建接收器是智慧型手機的一大突破。馬斯克（Elon Musk）去年曾說2023年定要開發一款手機可以直接連結星鏈，但迄今未實現，華為反而搶先做到。",
				"確實，在中國經濟動蕩之際，華為推出的七奈米手機給許多中國民眾打了一劑愛國主義強心劑。在上海及深圳等大城華為實體門市，排起了人龍搶購Mate 60 Pro手機。要價近7000元人民幣（約955美元）的新機銷售一空。",
				"中國官媒《環球時報》在8月底便在社論稱，華為在雷蒙多訪問期間開賣，被很多中國網民賦予了「在美國打壓之下昂起頭來」的更深層次含義，這些聲音應當被雷蒙多以及更多美國人聽到，並對華盛頓形成觸動力量。",
				"多家中國電商網站上出現背後印有雷蒙多頭像的新款手機外殼，微博則有網友後制的諷刺圖片，雷蒙多成為華為手機代言人，標題「我是雷蒙多，這次我為華為代言」。",
				"北京郵電大學經濟管理學院教授呂廷傑接受官媒央視採訪時就稱，從目前全球各機構對華為新機的拆解觀察，華為裝配的「麒麟9000s」七奈米晶片和其他該手機內的一萬多種零組件，基本上已實現了「國產化」。這意味著5G智能手機領域突破「卡脖子」問題。",
				"挑戰仍在",
				"賽格（Ahsel Sag）認為，中芯國際自去年以來一直在談論7奈米和5奈米技術，所以此次與華為合作推出配載7奈米晶片手機，並不是一個太大的驚喜。他告訴BBC中文，隨著現在半導體晶片，挑戰更精微的制程，需要依賴EUV進行光刻，但該機台目前無法銷售到中國，因此這個半導體面臨的挑戰只會變得更大。沒有它，中國將不得不發明新的光刻技術並研發自己的解決方案。",
				"劉佩真則強調，華為新機後續實際的成效則待觀察。她認為，中芯國際以DUV機台能衝刺的制程極限至多到七奈米制程，有關晶片的良率對一直晶片代工的一大挑戰，要大幅量產的成本也很高。",
				"許多分析師都強調，有關EUV光刻機機台的技術缺乏，仍然是中國半導體發展的一大軟肋。",
				"朱岳中就告訴BBC中文，在華為的新機中，由中芯代工製作的CPU系統，可能還是使用荷商艾司摩爾（ASML）生產的DUV光刻機，而非中國自製的光刻機。現在中國最先進的光刻機大概就屬上海微電子正在研發的28nm機台，雖然該機台預計2023年底出貨，但距離艾司摩爾生產的光刻機，至少有10年以上的技術差距。 雪上加霜的是，目前中芯在禁令下，也已買不到新的艾司摩爾光刻機了，現有的機台能撐多久是大挑戰。",
				"中國開發出5G或更高階的晶片「一直只是時間的問題，因為中國的科技實力本來就不弱。此次華為新機，中國自製率據估計約90%，但還是用了台灣大立光、穩懋，以及日本村田等公司的零組件」，朱岳中強調。 劉佩真也認為，高階的艾司摩爾EUV光刻機，在白宮要求下現在無法出售給中企，換言之，目前麒麟9000S 7奈米晶片，該是中芯國際在未來5至10年內能做到的最好制程技術。",
			},
		},
		{
			FileName: "002.html",
			Title:    "「一帶一路」倡議十週年 改變世界的成敗與未來",
			Author:   []string{"陳岩"},
			Category: "chinese-news",
			GUID:     "chinese-news-66762882",
			Tags:     []string{"中國", "政治", "經濟", "習近平"},
		},
		{
			FileName: "003.html",
			Title:    "中國經濟是一顆「定時炸彈」嗎？",
			Author:   []string{"尼克·馬什（Nick Marsh）"},
			Category: "business",
			GUID:     "business-66668559",
			Tags:     []string{"中國", "經濟", "習近平", "金融財經", "銀行業"},
		},
		{
			FileName: "004.html",
			Title:    "金正恩訪俄：朝鮮領導人進入俄羅斯國境，將與普京會面",
			Author:   []string{},
			Category: "world",
			GUID:     "world-66783048",
			Tags:     []string{"俄國", "弗拉基米爾·普京", "政治", "朝鮮", "烏克蘭局勢升溫", "軍事", "金正恩"},
			Content: []string{
				"朝鮮領導人金正恩已進入俄羅斯國境，預計將與普京總統舉行高峰會面。",
				"最新消息指，俄羅斯總統普京宣佈計劃前往東方太空發射場（Vostochny Cosmodrome），而金正恩預計也會前往該處。",
				"但是二人最終會面的地點，外界至今未有確認。",
				"克里姆林宮此前證實了金正恩的訪俄行程，但是莫斯科的聲明僅確認兩國領導人將會於遠東會面，未具體說明地點。",
				"韓國國防部較早前證實，金正恩的私人列車已於周二早上進入俄羅斯境內。",
				"朝鮮官方媒體發佈的照片顯示，金正恩在離開平壤之前，在他的裝甲列車上揮手。",
				"從啟程起，列車的行蹤就一直受到關注，外界預期金正恩進入俄羅斯後會需五至六個小時車程到達符拉迪沃斯托克。",
				"據朝中社報道，金正恩由高級政府官員陪同，其中包括軍方人員。",
				"但是，與此前外界預期不一樣的是，有報道指金正恩的專列在入境後似乎並未直接駛往符拉迪沃斯托克（中國舊稱「海參崴」），而是經過該城市後向北駛去。",
				"國營媒體俄羅斯新聞社（Ria）指，金正恩的裝甲列車駛入了符拉迪沃斯托克北邊的城市烏蘇里斯克，稍後又離開，駛向哈巴羅夫斯克（Khabarovsk，原名伯力）方向，其終點站仍然未明。",
				"該列車從邊境的哈桑火車站駛入俄羅斯，路透社引述匿名消息源指，金正恩在邊境曾下車與俄羅斯代表人員見面。",
				"分析人士指，如果目的地是符拉迪沃斯托克，列車應在到達烏蘇里斯克之後轉向南方。",
				"俄羅斯正在符拉迪沃斯托克主辦東方經濟論壇，外界此前預計金正恩與普京很可能在當地會面。",
				"會晤預計最早於當地時間周二進行——但克里姆林宮的聲明指，會面也有可能在「今後幾天」內進行。",
				"路透社報道，普京宣佈前往太空發射場時，沒有表明是否會在那裏與金正恩會面，僅表示：「等我到了那裏，你就會知道。」",
				"一名美國官員此前表示，由於俄羅斯正面臨烏克蘭的反攻，朝俄兩國很可能會討論軍火交易。",
				"據BBC的美國合作方媒體哥倫比亞廣播公司新聞台（CBS News）的報道，五角大樓預計朝俄領導人將會有「某種類型的會議」。",
				"如果金正恩與普京總統的峰會如期舉行，就將是朝鮮領導人四年多以來首次國際外訪，也是新冠疫情爆發以來的首次。",
				"美國官員在較早前向CBS表示，會面議程的一個重點是探討朝鮮是否可能向俄羅斯提供武器，支持其在烏克蘭的戰爭。",
				"金正恩上一次出國是在2019年，當時是朝鮮與時任美國總統特朗普（Donald Trump）的核裁軍談判破裂之後，金正恩前往符拉迪沃斯托克與普京總統舉行會面。",
				"有傳聞指，金正恩的火車包括至少20個防彈車廂，令它比一般的列車更重，速度最多不能超過59公里/小時（37英里/小時），前往符拉迪沃斯托克的旅程將需要整整一天。",
				"白宮此前表示，他們已獲得最新消息，俄羅斯和朝鮮之間的軍火談判正在「積極進展」。",
				"美國國家安全委員會發言人約翰·柯比爾（John Kirby）較早前表示，俄羅斯國防部長謝爾蓋·紹伊古（Sergei Shoigu）在最近一次訪朝期間曾試圖「說服平壤向俄羅斯出售火炮彈藥」。",
				"《紐約時報》較早前也曾報道，朝鮮可能尋求俄羅斯協助其發展太空項目。",
				"美國國務院在周一曾形容，俄羅斯向朝鮮這個「國際棄兒」伸出橄欖枝，是一次「戰略性失敗」。",
				"美國還一直持續警告朝鮮、中國及伊朗，不要援助俄羅斯在烏克蘭的戰事。",
				"克里姆林宮發言人德米特里·佩斯科夫（Dmitry Peskov）表示，俄羅斯對物華盛頓方面的警告「不感興趣」。",
				"克里姆林宮表示：「在落實我們與包括朝鮮在內的鄰國關係時，對我們來說重要的是兩國的利益，而不是來自華盛頓警告。",
				"」我們將會聚焦的是我們的兩個國家的利益。」",
				"總部位於華盛頓的卡內基國際和平研究院（Carnegie Endowment for International Peace）的安吉特·潘達（Ankit Panda）指出，俄羅斯和朝鮮各自有對方想要的東西。",
				"「現在重要的是雙方能否找到一個彼此都願意支付的價錢來換取對方的幫助，」他向BBC表示。",
				"俄羅斯可能會向朝鮮索要常規武器，包括炮彈和火箭炮彈藥，並以糧食和原材料作為交換，並在聯合國等國際論壇上繼續支持朝鮮。",
				"「這可能會開啟朝鮮向俄羅斯提供更精密武器的可能性，從而讓莫斯科能夠維持和補充自己的常規武器儲備。」",
				"俄羅斯被認為可能需要122毫米和152毫米炮彈，因為它的庫存已經告急，但是鑒於朝鮮的神秘，要確定其完整的火炮存量並不容易。",
				"金正恩與紹伊古在七月份會晤時，所展示的武器包括被認為是該國首款使用固體推進劑的火星系列（Hwasong）洲際彈道導彈。",
				"那也是自新冠疫情全球大流行以來，金正恩首次向外國客人打開國門。",
			},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.FileName,
			func(t *testing.T) {
				t.Parallel()
				f, err := os.Open(fmt.Sprintf("example/bbc/%s", tc.FileName))
				require.NoError(t, err)
				q := p.Parse(parser.NewTestQuery(200, "", f))
				require.NoError(t, q.Error)
				require.NotNil(t, q)

				require.Equal(t, tc.GUID, q.News.GUID)
				require.Equal(t, tc.Title, q.News.Title)
				require.Equal(t, tc.Category, q.News.Category)
				require.False(t, q.News.PubDate.IsZero())
				require.ElementsMatch(t, tc.Author, q.News.Author)
				require.ElementsMatch(t, tc.Tags, q.News.Tag)

				if len(tc.Content) > 0 {
					require.Equal(t, len(tc.Content), len(q.News.Content))
					for i, c := range q.News.Content {
						require.Equal(t, tc.Content[i], c)
					}
				}
			},
		)
	}
}

func TestNYTimes(t *testing.T) {
	p := parser.NewNYTimesParser()

	type testCase struct {
		FileName string
		Title    string
		Category string
		Author   []string
		GUID     string
		Tags     []string
		Content  []string
	}

	tcs := []testCase{
		{
			FileName: "001.html",
			Title:    "中國利用AI散播夏威夷山火陰謀論，對美虛假信息戰升溫",
			Category: "world",
			Author:   []string{"DAVID E. SANGER", "STEVEN LEE MYERS"},
			GUID:     "world-20230913-china-disinformation-ai",
			Tags:     []string{},
			Content: []string{
				"上個月，當毀滅性的山火迅猛橫掃毛伊島時，中國越來越機敏的信息鬥士們抓住了機會。",
				"他們忙著在互聯網上傳播虛假信息帖，號稱這場災難不是自然災害，而是美國進行祕密「氣象武器」實驗的結果。為了增強可信度，這些帖子還配有看來是用人工智慧程序生成的圖片，讓他們成為首批使用這種新工具為虛假信息運動增添真實性的人。",
				"2016年和2020年的美國總統大選期間，當俄羅斯開展駭客行動和虛假信息運動時，中國基本上置身事外。因此，對中國來說，把山火描述為美國情報部門和軍方的蓄意行為是一個戰術上的迅速轉變。",
				"中國此前的影響力運動一直集中於擴大宣傳，為其台灣政策和其他議題辯護。微軟和幾家不同機構的研究人員披露了該國的最新努力，表明中國政府正在更直接地嘗試在美國製造不和。",
				"就在中國發起這一努力之際，拜登政府和美國國會也在設法找到反擊中國、但不讓兩國陷入公開衝突的方式，並在設法降低人工智慧技術被用來放大虛假信息的風險。",
				"這場中國宣傳運動由來自微軟、Recorded Future、蘭德公司、NewsGuard，以及馬里蘭大學的研究人員發現，其影響尚且難以衡量，儘管初步跡象表明，幾乎沒有社群媒體用戶參與其中最離譜的陰謀論。",
				"微軟副董事長兼總裁布拉德·史密斯的研究人員對中國這場隱蔽的運動進行了分析，史密斯尖銳地批評了中國為自己的政治利益而利用自然災害的做法。",
				"「我認為這種做法對任何國家來說都不值得，更不用說渴望成為偉大國家的國家了，」史密斯週一接受採訪時說。",
				"利用毛伊島山火進行政治炒作的國家不只中國一個。俄羅斯也在網上發帖，強調美國在烏克蘭戰爭上花了多少錢，並聲稱這些錢本應該用在國內救災上。",
				"上述研究人員認為，中國正在建立一個可用於未來信息行動（包括美國明年的總統大選）的帳號網路。這正是俄羅斯在2016年美國總統大選前一年左右時間裡建立的模式。",
				"「這是在進入一個新方向，即放大與他們的某些利益——比如台灣——並無直接關係的陰謀論，」布萊恩·利斯頓說，他是總部位於麻薩諸塞州的網路安全公司Recorded Future的研究員。",
				"如果中國真要對美國明年的大選展開影響力行動的話，美國情報官員們近幾個月的評估是，中國很可能會試圖削弱拜登總統，同時提高前總統川普的形象。雖然對那些還記得川普用「中國病毒」的說法、試圖將大流行歸咎於北京的美國人來說，這似乎有悖常理，但情報官員已得出結論：中國領導人更喜歡川普。川普提出將美國撤出日本、韓國，以及亞洲其他地區，而拜登的做法包括切斷中國獲得最先進晶片及其生產設備的途徑。",
				"中國推出有關夏威夷火災的陰謀論之前，去年秋天，拜登曾在峇里島向中國國家主席習近平指責北京在傳播各種虛假信息上起的作用。據政府官員稱，拜登憤怒地批評習近平讓中國散布虛假指控，稱美國在烏克蘭運行生物武器實驗室。",
				"研究人員和政府官員表示，沒有跡象表明俄羅斯和中國在信息行動上合作，但兩國經常附和對方的信息，尤其是在批評美國政策時。兩國的共同努力意味著，虛假信息戰的新階段即將開始，而人工智慧工具的使用將使這場戰爭更加激烈。",
				"「雖然我們沒有直接證據表明中國和俄羅斯在這些行動上進行協調，但我們確實看到了一致性和某種同步性，」蘭德公司研究員威廉·馬切利諾說，他寫了一份新報告，警告人工智慧將讓全球影響力行動出現「重要飛躍」。",
				"夏威夷的山火和如今發生的許多自然災害一樣，幾乎從一開始就催生了許多謠言、虛假報導和陰謀論。",
				"馬里蘭大學的情報與安全應用研究實驗室的研究員卡羅琳·艾米·奧爾·布埃諾報告稱，8月9日，也就是山火發生的第二天，俄羅斯在Twitter（現名X的社群媒體平台）上發起協調行動。",
				"俄羅斯用一個沒有幾名關注者的不起眼帳號傳播「關注夏威夷、而非烏克蘭」(Hawaii, not Ukraine)的說法，旨在削弱美國向烏克蘭提供的軍事援助。它通過美國的布萊巴特新聞網等一系列保守派或右翼帳號，最終通過俄羅斯官方媒體，傳播給了成千上萬名用戶。",
				"中國官媒經常附和俄羅斯的主題，尤其是對美國的敵意。但這次，中國展開了一場不同類型的虛假信息運動。",
				"據Recorded Future最先發布的報告，中國政府發起了一場隱蔽行動，將山火歸咎於「氣象武器」。該公司在8月中旬找到了大量帖子，錯誤地聲稱英國的外國情報機構軍情六處爆料「山火背後的驚人真相」。使用完全一樣語言的帖子也出現在互聯網的各種社群媒體平台上，包括Pinterest、Tumblr、Medium，以及藝術家使用的日本網站Pixiv。",
				"其他傳播不實信息的帳號也發了類似的內容，還經常配有貼錯標籤的影片，其中一個名為The Paranormal Chic的熱門TikTok帳號發的影片顯示的是智利的一起變壓器爆炸事故。據Recorded Future的報告，中國的這些內容經常呼應並放大美國的陰謀論者、以及包括白人至上主義者在內的極端分子的帖子。",
				"中國這場宣傳運動橫跨多個主要社群媒體平台，使用多種語言，表明其目的是引起全球受眾的注意。微軟的威脅分析中心找到了用31種語言發布的不實帖子，包括法語、德語、義大利語，以及伊博語、奧迪亞語和瓜拉尼語等使用者更少的帖子。",
				"被微軟研究人員認定為人工生成的夏威夷山火圖片出現在多個平台上，包括一個用荷蘭語發在Reddit上的帖子。「這些用人工智慧生成的具體圖片似乎只被（此次運動中用的中國帳號）使用。它們似乎沒有出現在網上其他地方，」微軟在一份報告中說。",
				"微軟的威脅分析中心總經理克林特·瓦茨表示，中國似乎在使用俄羅斯的影響力戰術，為影響美國和其他國家的政治奠定基礎。",
				"「這像是俄羅斯2015年的做法，」他說，他指的是在2016年美國總統大選期間在網上開展大規模影響力運動之前，俄羅斯建立機器人帳號和虛假帳號的做法。「如果我們看一下其他行動者是如何做的，就會發現，他們正在打造能力。現在他們正在創建隱蔽帳號。」",
				"自然災害常常是虛假信息運動的重點，因為它能讓不懷好意的行為者利用人們的情緒，指責政府對災害準備不足、或應對乏力。這樣做的目的可能是破壞人們對具體政策的信心，例如美國對烏克蘭的支持；或者更廣泛地煽動內部不和。通過暗示美國正在本國公民身上測試或使用祕密武器，中國的這一努力似乎也是為了將美國描繪成不計後果的軍國主義大國。",
				"「我們總是能在人道主義災難發生後團結起來，為遭受地震、颶風或火災破壞的人民提供幫助，」史密斯週二向美國國會介紹微軟的一些發現時說。「現在看到這種利用災難的事情，令人深感不安，我認為國際社會應該劃定紅線、禁止這種行為。」",
			},
		},
		{
			FileName: "002.html",
			Title:    "中國充滿風險，為何跨國企業仍難以離開",
			Author:   []string{"艾莎"},
			Category: "business",
			GUID:     "business-20230911-china-us-business",
			Tags:     []string{},
			Content: []string{
				"過去數十年裡，美國的企業老總們將中國視作一棵搖錢樹。他們誇張地稱讚中國的數億消費者，把這個市場稱為「最大的機會之一」，並預測本世紀將是「中國的世紀」。",
				"現在，這些高管們在最近訪華後帶回了更清醒的看法。在中國做生意的西方公司正面臨著幾年前難以想像的壓力。中國經濟面臨重重困難，與美國的關係陷入緊張。持續三年的出入境限制以及商業活動事實上的停止造成了尚未癒合的裂痕。",
				"中國取消「新冠清零」政策、重新開放已經九個月了，企業正在努力設法面對一個嚴峻的現實：雖然總值18萬億美元的中國經濟充滿了危險，但仍不可忽視，也難以離開。撤離可能意味著在未來的全球競爭中失去優勢。許多西方公司仍將中國業務視為它們押下的長期賭注，儘管回報受到了風險的牽累。",
				"「首席執行官們已經意識到，他們需要降低一些風險，」大成全球顧問公司的高級顧問薄邁倫(Myron Brilliant)說。「他們不想忽視中國市場，但對目前的環境都有充分的認識。」",
				"讓人擔心的事情很多。警方突襲西方公司駐華辦事處、政府的高額罰款、破壞收購交易、限制數據傳輸的法規，以及涉及範圍廣泛的反間諜法，都增加了做生意的成本。還有被稱為灰天鵝的其他風險，即罕見但並非不可想像的事件，例如又一場大流行病、更多的經濟制裁，或公開的跨境衝突。這些擔憂加在一起，是一種被美國商務部長雷蒙多最近描述為美國企業中存在的中國「不適合投資」的感覺。",
				"餘波能很快地感受到。本週，有關中國政府將禁止政府機構和其他國控實體的員工使用iPhone的報導出來後，蘋果公司的股價下跌了6%，導致其市值蒸發了近2000億美元。",
				"中國不斷惡化的經濟前景加劇了企業的擔憂，讓它們更難有理由在中國進行更多投資。曾經三年無法入境中國的外國企業高管們終於能開始前往中國，與那裡的員工見面。許多人曾預計會看到中國經濟的強勁復甦。",
				"然而，一些高管回到美國後開始擔心，中國官員對自己應對經濟衰退的能力過於自信。私下裡，外國企業的領導人已對中國企業不再投資國內業務感到震驚。他們問，如果中國自己的民營企業對經濟都沒有信心的話，我們為什麼應該在中國投資呢？",
				"「企業董事會會議上有關中國的討論正在不可避免地轉向更加謹慎，」華盛頓戰略與國際研究中心的中國問題專家白明(Jude Blanchette)說。原因在於經濟放緩，他說，以及「中國政府難以琢磨的懲罰性監管行為和走向極權主義的趨勢，再就是美國政府將技術和投資引導到其他市場的行動」。",
				"美國官員的立場也讓事情變得更加複雜，他們的態度已轉向與中國對立。在中國採取一切照舊的做法可能意味著被美國立法者傳喚到聽證會上。「如果你對中國有任何正面說法的話，立法者們會讓你在那裡坐得很不舒服，」發動機跨國公司康明斯的發言人喬恩·米爾斯說，這是一家有百年歷史的美國公司。",
				"在聽證會上受盤問帶來名譽和法律後果。由威斯康辛州共和黨眾議員邁克·加拉格爾擔任主席的與中國競爭特別委員會有傳喚證人出席聽證會的權力和政治影響力。該委員會並非呼籲停止與中國建立夥伴關係的唯一聲音。",
				"維吉尼亞州的共和黨人州長格倫·楊金把福特汽車與一家中國公司達成的協議描述為中共的「特洛伊木馬」，福特將使用這家中國公司的電池技術在密西根州建一家電動汽車電池廠。楊金已阻止福特把工廠建在維吉尼亞州。",
				"莫德納在中國研究、開發和生產使用mRNA藥物的決定被佛羅里達州的共和黨參議員馬可·盧比奧說成是「對美國納稅人的背叛，是他們的辛勤勞動所得讓這項技術成為可能」。",
				"特斯拉在上海建一個大型電池工廠的計劃已引發了加拉格爾對特斯拉過於依賴「中國市場准入」的質疑。",
				"美國企業極力在政治審查與這樣一種信念之間尋找平衡，也就是如果不在研究與創新上與中國的公司展開競爭合作的話，就有落後的危險，因為中國的競爭對手們將在全球市場上擊敗它們。",
				"福特最近與中國寧德時代新能源科技有限公司建立合作夥伴關係，不是以增加福特在中國的業務的方式（以免冒下在美國國內受批評的風險），而是以福特在密西根州建一家獨資且自己運營的電池廠的方式。福特稱合作項目將在美國創造2500個就業崗位。這家耗資35億美元的工廠將使用寧德時代的技術製造電池（寧德時代是全球最大的電動汽車電池製造商），以「幫助我們更快地生產更多的電動汽車」，福特公司的執行董事長小威廉·克萊·福特說。",
				"儘管如此，共和黨議員們表示，正在調查該合作項目，因為擔心寧德時代與新疆強迫勞動有關，聯合國指出中國在西部新疆地區存在系統性侵犯人權的行為。",
				"就製藥而言，中國已明確表示，外國企業需要改變傳統的經營方式，不要只將在國外研發的藥物帶進中國市場，而是要與本土科學家合作，投資藥物研發。",
				"對莫德納來說，中國龐大的患者群體、藥物研發方面的大量投資，以及臨床試驗資源，可能是它決定與中國合作的部分原因。有報導稱，莫德納曾在數月前拒絕了中國讓其交出新冠疫苗背後的知識產權的要求。莫德納正面臨著新冠疫苗需求下降的問題，該疫苗是公司在商業上可行的唯一產品，進入中國將讓莫德納在全球最大的製藥業市場之一開發其他使用mRNA技術的疫苗。",
				"在習近平擔任最高領導人的十年裡，中國政府已在他的控制下將注意力大幅轉向內部。「從結構上看，目前的政府定位與前幾屆政府的有很大不同，」艾意凱諮詢董事合伙人陳瑋說。「非常強調中國的崛起，但這對西方公司意味著什麼呢？」",
				"即使高管們希望像一些美國立法者正在推動的那樣與中國脫鉤，許多公司都表示這是不合理的。康明斯的米爾斯說，切割中國業務不可行。這家製造發動機、發電機，以及汽車零部件的企業在中國擁有21家工廠，約20%的利潤來自中國。",
				"「我們在中國的成功已帶來了我們在全球的成功，以及美國的就業增長，」他補充道。",
				"其他公司也有同感。",
				"「我認為，讓美國人民了解美國與中國的關係很重要，我們需要找到一種相處的方式，」RTX首席執行官格雷格·海耶斯在今年早些時候接受CNBC採訪時說。RTX是一家航空航天和國防項目承包商，前身是雷神公司，目前在中國擁有兩家生產商用發動機、航空系統和客艙的子公司。海耶斯說，將供應鏈撤出中國不切實際，中國市場「對美國經濟來說太大、太重要、太有必要」。",
				"但激烈的競爭以及不斷增加的地緣政治、戰略和財務成本已降低了美國企業界對中國曾經有過的熱情。",
				"隨著中國面臨其幾十年來最大的經濟挑戰，許多跨國公司正在世界其他地方尋求增長機會，大成全球顧問公司的薄邁倫說道，他曾任美國全國商會副會長。",
				"「由於中國經濟走向存在一定程度的不確定性，企業高管們照常行事就是玩忽職守，」他說。",
			},
		},
		{
			FileName: "003.html",
			Title:    "蘋果發布iPhone 15系列，轉用USB-C接口",
			Author:   []string{"TRIPP MICKLE"},
			Category: "technology",
			GUID:     "technology-20230913-apple-iphone-15-usb-c",
			Tags:     []string{},
		},
		{
			FileName: "004.html",
			Title:    "馬斯克為何對「X」情有獨鍾",
			Author:   []string{"STELLA BUGBEE"},
			Category: "style",
			GUID:     "style-20230727-gen-x-elon-musk-new-logo",
			Tags:     []string{},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.FileName,
			func(t *testing.T) {
				t.Parallel()
				f, err := os.Open(fmt.Sprintf("example/nytimes/%s", tc.FileName))
				require.NoError(t, err)
				q := p.Parse(parser.NewTestQuery(200, "", f))
				require.NoError(t, q.Error)
				require.NotNil(t, q)
				require.Equal(t, tc.GUID, q.News.GUID)
				require.Equal(t, tc.Title, q.News.Title)
				require.Equal(t, tc.Category, q.News.Category)
				require.False(t, q.News.PubDate.IsZero())
				require.ElementsMatch(t, tc.Author, q.News.Author)
				require.ElementsMatch(t, tc.Tags, q.News.Tag)

				if len(tc.Content) > 0 {
					require.Equal(t, len(tc.Content), len(q.News.Content))
					for i, c := range q.News.Content {
						require.Equal(t, tc.Content[i], c)
					}
				}
			},
		)
	}
}

type urlMatcher struct {
	*url.URL
}

func (m urlMatcher) Matches(x any) bool {
	u, ok := x.(*url.URL)
	if !ok {
		return false
	}
	return m.URL.String() == u.String()
}

func TestPraserRepo(t *testing.T) {
	ctl := gomock.NewController(t)
	p := mock_parser.NewMockParser(ctl)

	domains := []string{
		"localhost:80",
		"localhost:8080",
		"localhost:443",
	}
	p.EXPECT().
		Domain().
		Times(2).
		Return(domains)

	ru := "http://localhost:80/static/001.html"
	u, err := url.Parse(ru)
	require.NoError(t, err)

	p.EXPECT().
		ToGUID(urlMatcher{u}).
		Times(2).
		Return(u.Path)

	now := time.Now()
	news := &parser.News{
		Title:       "title",
		Link:        u,
		Description: "description",
		Language:    "en-us",
		Author:      []string{"Charles", "Moby"},
		Category:    "technology",
		GUID:        u.Path,
		PubDate:     now,
		Content:     []string{"first paragraph", "second paragraph", "third paragraph"},
		Tag:         []string{"tag1", "tag2", "tag3"},
		RelatedGUID: []string{"/static/002.html", "/static/003.html", "/static/004.html"},
	}

	repo := parser.NewParserRepo(p)
	require.True(t, repo.Has("localhost:80"))
	require.True(t, repo.Has("localhost:8080"))
	require.True(t, repo.Has("localhost:443"))
	require.False(t, repo.Has("host-not-found"))

	require.Equal(t, u.Path, repo.ToGUID(u))
	require.ElementsMatch(t, domains, repo.Domain())

	p80, ok := repo[u.Host]
	require.True(t, ok)
	require.NotNil(t, p80)

	require.Equal(t, domains, p80.Domain())
	require.Equal(t, u.Path, p80.ToGUID(u))

	q := parser.NewQuery(ru)
	require.NotNil(t, q)

	p.EXPECT().
		Parse(gomock.Any()).
		Times(1).
		DoAndReturn(func(q *parser.Query) *parser.Query {
			jsn, _ := json.Marshal(news)

			n := &parser.News{}
			_ = json.Unmarshal(jsn, n)
			q.News = n
			return q
		})

	q = repo.Parse(q)
	require.NotNil(t, q)
	require.Equal(t, q.News.Title, news.Title)
	require.Equal(t, q.News.Link, news.Link)
	require.Equal(t, q.News.Description, news.Description)
	require.Equal(t, q.News.Language, news.Language)
	require.Equal(t, q.News.Author, news.Author)
	require.Equal(t, q.News.Category, news.Category)
	require.Equal(t, q.News.GUID, news.GUID)
	require.Equal(t, q.News.Content, news.Content)
	require.Equal(t, q.News.Tag, news.Tag)
	require.Equal(t, q.News.RelatedGUID, news.RelatedGUID)
	require.Greater(t, 10*time.Second, q.News.PubDate.Sub(news.PubDate))

	u, err = url.Parse("http://not-found-in-repo")
	require.NoError(t, err)

	q = parser.NewQuery(u.String())
	q = repo.Parse(q)
	require.NotNil(t, q)
	require.Empty(t, q.News)
	require.Equal(t, parser.ErrParserNotFound, q.Error)
	require.Equal(t, u.Path, repo.ToGUID(u))

	q = &parser.Query{RawURL: u.String()}
	q = repo.Parse(q)
	require.NotNil(t, q)
	require.Empty(t, q.News)
	require.Equal(t, parser.ErrParserNotFound, q.Error)
	require.Equal(t, u.Path, repo.ToGUID(u))

	q = &parser.Query{RawURL: ":not-url"}
	q = repo.Parse(q)
	require.NotNil(t, q)
	require.Error(t, q.Error)
}

func TestDefaultParser(t *testing.T) {
	repo := parser.GetDefaultParser()
	require.NotNil(t, repo)

	for _, d := range repo.Domain() {
		require.True(t, parser.Has(d))
	}

	var q *parser.Query
	ru := "https://unknown.host/news/very-important-article"
	u, err := url.Parse(ru)
	require.NoError(t, err)

	require.False(t, parser.Has(u.Host))

	q = parser.ParseRawURL(ru)

	require.ErrorIs(t, q.Error, parser.ErrParserNotFound)
	require.Nil(t, q.News)

	q = parser.ParseURL(u)
	require.ErrorIs(t, q.Error, parser.ErrParserNotFound)
	require.Nil(t, q.News)

	q = parser.Parse(&parser.Query{RawURL: ru})
	require.Nil(t, q.News)
	require.ErrorIs(t, q.Error, parser.ErrParserNotFound)

	q = parser.Parse(&parser.Query{URL: u})
	require.Nil(t, q.News)
	require.ErrorIs(t, q.Error, parser.ErrParserNotFound)
}

func TestParserRepoForDifferentDomain(t *testing.T) {
	type testCase struct {
		Name     string
		RawURL   string
		URL      *url.URL
		Parser   parser.Parser
		FileName string
		Title    string
		GUID     string
	}

	tcs := []testCase{
		{
			Name:     "UDN",
			RawURL:   "https://udn.com/news/story/123707/7398504",
			Parser:   parser.NewUDNParser(),
			FileName: "/udn/news/001.html",
			Title:    "核汙水排入海 中日緊張升溫",
			GUID:     "123707-7398504",
		},
		{
			Name:     "PTS",
			RawURL:   "https://news.pts.org.tw/article/652596",
			Parser:   parser.NewPTSParser(),
			FileName: "/pts/001.html",
			Title:    "美韓軍演首納太空軍 一文了解太空軍是什麼",
			GUID:     "652596",
		},
		{
			Name:     "CNA",
			RawURL:   "https://www.cna.com.tw/news/ahel/202308230179.aspx",
			Parser:   parser.NewCNAParser(),
			FileName: "/cna/001.html",
			Title:    "新北攤商涉用豬肉混充羊肉 最重可罰400萬加坐牢",
			GUID:     "ahel-202308230179",
		},
		{
			Name:     "ETtoday",
			RawURL:   "https://www.ettoday.net/news/20230904/2575420.htm",
			Parser:   parser.NewEtTodayParser(),
			Title:    "快訊／侯友宜將有「重大宣布」！　她喊：大家拭目以待",
			FileName: "/ettoday/001.html",
			GUID:     "20230904-2575420",
		},
		{
			Name:     "BBC-ZH",
			RawURL:   "https://www.bbc.com/zhongwen/trad/chinese-news-66748432",
			Parser:   parser.NewBBCParser(),
			FileName: "/bbc/001.html",
			Title:    "華為：新款5G手機橫空出世 中國半導體突破美國封鎖？",
			GUID:     "chinese-news-66748432",
		},
	}

	repo := parser.ParserRepo{}
	for i := range tcs {
		u, err := url.Parse(tcs[i].RawURL)
		require.NoError(t, err)
		tcs[i].URL = u
		repo[u.Host] = tcs[i].Parser
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("Case %d: %s", i+1, tcs[i].Name),
			func(t *testing.T) {
				require.True(t, repo.Has(tc.URL.Host))

				f, err := os.Open("example/" + tc.FileName)
				require.NoError(t, err)

				q := parser.NewTestQuery(200, "", f)
				q.URL = tc.URL
				q.RawURL = tc.RawURL
				q = repo.Parse(q)
				require.NoError(t, q.Error)
				require.Equal(t, tc.Title, q.News.Title)
				require.Equal(t, tc.GUID, q.News.GUID)
			},
		)
	}
}
