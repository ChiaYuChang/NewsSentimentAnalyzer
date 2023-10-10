package parser_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/parser"
	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

func TestJsonLDUnmarshal(t *testing.T) {
	var err error
	var s string = `{
            "@context": "https://schema.org",
            "@type": "NewsArticle",
            "articleSection": "technology",
			"about": [
                {
                    "@type": "Thing",
                    "name": "bussiness",
                    "sameAs": [
                        "http://example.com/bussiness"
                    ]
                },
				{
                    "@type": "Person",
                    "name": "Keven",
                    "sameAs": [
                        "http://example.com/keven"
                    ]
                }
			],
            "headline": "page headline",
            "author": {
                "@type": "Person",
                "name": "John"
            },
            "url": "https://test.page/001.html",
            "description": "page description",
            "keywords": "kw1,kw2,kw3,kw4",
            "datePublished": "2023-08-25T19:04:00+08:00",
            "dateModified": "2023-08-25T19:04:00+08:00"
    }`
	if s[0] == '{' && s[len(s)-1] == '}' {
		s = fmt.Sprintf("[%s]", s)
	}

	ldjson1 := parser.JsonLDList{}
	err = json.Unmarshal([]byte(s), &ldjson1)
	require.NoError(t, err)
	require.Equal(t, "NewsArticle", ldjson1.GetByType("NewsArticle").Type)
	require.Equal(t, "technology", ldjson1.GetByType("NewsArticle").ArticleSection)
	require.Equal(t, "page headline", ldjson1.GetByType("NewsArticle").Headline)
	require.Equal(t, "John", ldjson1.GetByType("NewsArticle").Author[0].Name)
	require.Equal(t, "bussiness", ldjson1.GetByIndex(0).About[0].Name)
	require.Equal(t, "Keven", ldjson1.GetByIndex(0).About[1].Name)
	require.Equal(t, "https://test.page/001.html", ldjson1.GetByType("NewsArticle").URL.String())
	require.ElementsMatch(t, []string{"kw1", "kw2", "kw3", "kw4"}, ldjson1.GetByType("NewsArticle").Keywords)

	data, err := json.MarshalIndent(ldjson1, "", "\t")
	require.NoError(t, err)

	ldjson2 := parser.JsonLDList{}
	err = json.Unmarshal(data, &ldjson2)
	require.NoError(t, err)
	require.NoError(t, err)
	require.Equal(t, "NewsArticle", ldjson2.GetByIndex(0).Type)
	require.Equal(t, "technology", ldjson2.GetByIndex(0).ArticleSection)
	require.Equal(t, "page headline", ldjson2.GetByIndex(0).Headline)
	require.Equal(t, "John", ldjson2.GetByIndex(0).Author[0].Name)
	require.Equal(t, "bussiness", ldjson2.GetByIndex(0).About[0].Name)
	require.Equal(t, "Keven", ldjson2.GetByIndex(0).About[1].Name)
	require.Equal(t, "https://test.page/001.html", ldjson2.GetByIndex(0).URL.String())
}

func TestCondenseJsonLDUnmarshal(t *testing.T) {
	data := []byte(`{
        "@context": "https://schema.org",
        "@type": "NewsArticle",
        "name": "美韓軍演首納太空軍 一文了解太空軍是什麼 ｜ 公視新聞網 PNN",
        "description": "太空軍是美國軍隊6大軍種的最新分支，在美韓近期的例行軍演中，首度加入演習。其實美國太空軍在先期的任務規劃未包括直接作戰，一文了解美國太空軍在做什麼。",
        "image": [
                "https://news-data.pts.org.tw/media/170949/conversions/cover-webp.webp",
                "https://news-data.pts.org.tw/media/170949/conversions/cover-webp.webp"
        ],
        "author": {
                "@type": "person",
                "name": "陳宥蓁／綜合報導"
        },
        "articleBody": "美國和南韓於昨（21）日起為期10天舉行「乙支自由之盾」年度聯合軍演，為防堵北韓網路入侵，還首度將美軍成立不久的「太空軍」加入演習。駐韓美軍先前指出，太空軍的加入會使作戰更具多樣性。美韓例行軍演展開 首加入太空軍防北韓網路入侵美國太空軍是什麼？美國軍隊共分為6大軍種，太空軍於2019年12月成立，屬最新分支，總編制約為1萬6千人，主要工作是組織、訓練和裝備太空部隊，不過在先期任務規劃還未包括直接作戰，為的是保護美國在太空的資產及利益，阻止他人藉由太空進行侵略。在太空軍成立之際，美國前總統川普曾說，「太空是世界上最新的戰鬥領域。」而今年年初，太空軍作戰部長薩茲曼便指出，中國和俄羅斯已發展出反衛星飛彈和軌道攔截能力等技術，是美國在太空領域的兩大威脅。2020年5月於白宮舉行的美國太空軍軍旗亮相儀式。（圖／美聯社）美國太空軍如何分工？太空軍作戰部長辦公室（OCSO）是太空軍的最高指揮機構，其下有太空軍戰地司令部（Space Force Field Command），3個司令部又有各自的三角洲部隊（Delta），並有中隊（Squadron）提供支援。太空作戰司令部（SpOC）負責做好戰鬥準備，監視與偵察全球情報，並與外界夥伴合作。太空系統司令部（SSC）負責為開發、購買和部署太空武裝系統，監督國防部衛星和其他空間系統的發射操作、維護等，與太空相關的研究工作。太空訓練和準備司令部（STARCOM）負責培訓和教育人員，並制定太空術語和戰術等。新單位用「死神」當隊徽？美國太空軍近期成立了新單位ISRS，是首支恐會摧毀敵方衛星的部隊，與採用「死神」當隊徽的意涵相互呼應；而死神眼中的北極星芒，則代表安全的導引。ISRS負責的任務是阻斷對美國使用衛星系統造成的威脅，包含地面上可能造成影響的雷射、干擾訊號裝置，或駭入衛星系統的資安攻擊等。他們將分析潛在目標，並進行識別和追蹤。今年8月11日於ISRS部隊啟用儀式上揭曉的隊徽。（圖／USSF）美國境外也有太空軍？美國太空軍去（2022）年11月在夏威夷珍珠港的基地，首度在本土以外設立區域指揮中心，除了控管印太地區的太空軍事行動，也是為了防禦中國威脅，並包括俄羅斯及北韓。而同年12月，駐韓美軍也成立了美國太空軍的海外太空部隊，負責監控和追蹤北韓及鄰近地區飛彈，還有全球定位、衛星通訊等任務。",
        "url": "https://parser.pts.org.tw/article/652596",
        "headline": "美韓軍演首納太空軍 一文了解太空軍是什麼 ｜ 公視新聞網 PNN",
        "datePublished": "2023-08-22T10:28:47.000000Z",
        "keywords": [
                "太空",
                "空軍",
                "美國",
                "衛星",
                "部隊",
                "駐韓美軍"
        ],
        "publisher": {
                "@type": "Organization",
                "name": "公視新聞網",
                "url": "https://parser.pts.org.tw/",
                "logo": {
                        "@type": "ImageObject",
                        "url": "https://d3prffu8f9hpuw.cloudfront.net/news/static/public/images/logo.png"
                }
        },
        "mainEntityOfPage": "https://parser.pts.org.tw/article/652596",
        "dateModified": "2023-08-22T10:29:16.000000Z"
	}`)

	var ldjson parser.JsonLD

	require.True(t, json.Valid(data))

	err := json.Unmarshal(data, &ldjson)
	require.NoError(t, err)
}

func TestJsonLDList(t *testing.T) {
	data := []byte(`{
		"@context": "http://schema.org",
		"@graph": [
			{
				"@type": "NewsArticle",
				"url": "https://www.bbc.com/zhongwen/trad/chinese-news-66748432",
				"publisher": {
					"@type": "NewsMediaOrganization",
					"name": "BBC News 中文",
					"publishingPrinciples": "https://www.bbc.com/zhongwen/trad/institutional-51359584",
					"logo": {
						"@type": "ImageObject",
						"width": 1024,
						"height": 576,
						"url": "https://parser.files.bbci.co.uk/ws/img/logos/og/zhongwen.png"
					}
				},
				"thumbnailUrl": "https://ichef.bbci.co.uk/news/1024/branded_zhongwen/150AF/production/_131019168_gettyimages-1647674967.jpg",
				"image": {
					"@type": "ImageObject",
					"width": 1024,
					"height": 576,
					"url": "https://ichef.bbci.co.uk/news/1024/branded_zhongwen/150AF/production/_131019168_gettyimages-1647674967.jpg"
				},
				"mainEntityOfPage": {
					"@type": "WebPage",
					"@id": "https://www.bbc.com/zhongwen/trad/chinese-news-66748432",
					"name": "華為：新款5G手機橫空出世 中國半導體突破美國封鎖？"
				},
				"headline": "華為：新款5G手機橫空出世 中國半導體突破美國封鎖？",
				"description": "在美國商務部長雷蒙多結束北京行前夕，中國電子巨擘華為無預警發售此款手機，似乎在告訴雷蒙多，中國半導體不會被美國封鎖打垮。",
				"datePublished": "2023-09-12T07:56:32.000Z",
				"dateModified": "2023-09-12T07:56:32.000Z",
				"inLanguage": {
					"@type": "Language",
					"name": "Chinese",
					"alternateName": "zh-hant"
				},
				"about": [
					{
						"@type": "Thing",
						"name": "技術新知",
						"sameAs": [
							"http://dbpedia.org/resource/Technology"
						]
					},
					{
						"@type": "Thing",
						"name": "中美關係",
						"sameAs": [
							"http://dbpedia.org/resource/China%E2%80%93United_States_relations"
						]
					},
					{
						"@type": "Person",
						"name": "習近平",
						"sameAs": [
							"http://dbpedia.org/resource/Xi_Jinping"
						]
					},
					{
						"@type": "Place",
						"name": "中國",
						"sameAs": [
							"http://dbpedia.org/resource/China"
						]
					},
					{
						"@type": "Thing",
						"name": "華為",
						"sameAs": [
							"http://dbpedia.org/resource/Huawei"
						]
					},
					{
						"@type": "Thing",
						"name": "經濟",
						"sameAs": [
							"http://dbpedia.org/resource/Economy"
						]
					},
					{
						"@type": "Thing",
						"name": "貿易",
						"sameAs": [
							"http://dbpedia.org/resource/Trade"
						]
					}
				],
				"author": {
					"@type": "NewsMediaOrganization",
					"name": "BBC News 中文",
					"logo": {
						"@type": "ImageObject",
						"width": 1024,
						"height": 576,
						"url": "https://parser.files.bbci.co.uk/ws/img/logos/og/zhongwen.png"
					},
					"noBylinesPolicy": "https://www.bbc.com/zhongwen/trad/institutional-51359584#authorexpertise"
				}
			}
		]
	}`)

	re := regexp.MustCompile(`"@graph": [(.+)]`)
	dd := re.FindSubmatch(data)
	for i, d := range dd {
		t.Logf("%d: %q\n", i, string(d))
	}
	// t.Log(string(data))

	require.True(t, json.Valid(data))

	var ldjson parser.JsonLD
	err := json.Unmarshal(data, &ldjson)
	require.NoError(t, err)
	t.Log(ldjson)
}

func genRdmJsonLD() parser.JsonLD {
	u, _ := url.Parse(fmt.Sprintf(
		"https://%s.%s.%s/%s.html",
		rg.Must(rg.AlphabetLower.GenRdmString(5)),
		rg.Must(rg.AlphabetLower.GenRdmString(3)),
		rg.Must(rg.AlphabetLower.GenRdmString(3)),
		rg.Must(rg.AlphaNum.GenRdmString(10)),
	))

	keywords := make([]string, 10)
	for i := 0; i < len(keywords); i++ {
		keywords[i] = rg.Must(rg.AlphaNum.GenRdmString(5))
	}

	rt := rg.GenRdnTimes(2, time.Now().Add(-24*30*time.Hour), time.Now())
	jld := parser.JsonLD{
		Type:           rg.Must(rg.AlphaNum.GenRdmString(5)),
		Headline:       rg.Must(rg.AlphaNum.GenRdmString(10)),
		Description:    rg.Must(rg.AlphaNum.GenRdmString(10)),
		URL:            u,
		ArticleSection: rg.Must(rg.AlphabetLower.GenRdmString(8)),
		Author: []parser.JsonLDObject{
			{
				Type: "@person",
				Name: rg.Must(rg.Alphabet.GenRdmString(10)),
			},
		},
		Keywords:    parser.CSL(keywords),
		PublishedAt: rt[0],
		UpdatedAt:   rt[1],
	}
	return jld
}

func TestParseJsonLD(t *testing.T) {
	type testCase struct {
		N          int
		JsonLDObj  map[string]parser.JsonLD
		JsonLdData []byte
	}

	tcs := []testCase{
		{N: 1, JsonLDObj: make(map[string]parser.JsonLD)},
		{N: 3, JsonLDObj: make(map[string]parser.JsonLD)},
		{N: 5, JsonLDObj: make(map[string]parser.JsonLD)},
	}

	for i := range tcs {
		jlds := make([]parser.JsonLD, tcs[i].N)
		for j := range jlds {
			jld := genRdmJsonLD()
			tcs[i].JsonLDObj[jld.Type] = jld
			jlds[j] = jld
		}
		data, err := json.MarshalIndent(jlds, "", "\t")
		require.NoError(t, err)
		tcs[i].JsonLdData = data
	}

	allJlds := parser.JsonLDList{
		JsonLD:      make([]*parser.JsonLD, 0),
		TypeToIndex: make(map[string]int),
	}
	n := 0
	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("Case %d", i+1),
			func(t *testing.T) {
				var jlds1 parser.JsonLDList
				err := json.Unmarshal(tc.JsonLdData, &jlds1)
				require.NoError(t, err)

				for jtype := range jlds1.TypeToIndex {
					require.Equal(t, tc.JsonLDObj[jtype].Headline, jlds1.GetByType(jtype).Headline)
					require.Equal(t, tc.JsonLDObj[jtype].Description, jlds1.GetByType(jtype).Description)
					require.Equal(t, tc.JsonLDObj[jtype].URL, jlds1.GetByType(jtype).URL)
				}

				data, err := json.MarshalIndent(jlds1, "", "    ")
				require.NoError(t, err)

				var jlds2 parser.JsonLDList
				err = json.Unmarshal(data, &jlds2)
				require.NoError(t, err)

				for jtype := range jlds1.TypeToIndex {
					require.Equal(t, tc.JsonLDObj[jtype].Headline, jlds2.GetByType(jtype).Headline)
					require.Equal(t, tc.JsonLDObj[jtype].Description, jlds2.GetByType(jtype).Description)
					require.Equal(t, tc.JsonLDObj[jtype].URL, jlds2.GetByType(jtype).URL)
				}
				require.Equal(t, jlds1.String(), jlds2.String())
				allJlds.Merge(jlds1)
				n += tc.N
			},
		)
	}
	require.Equal(t, n, allJlds.Len())
}

func TestSelectorBuilder(t *testing.T) {
	htmldoc := `<!DOCTYPE html>
	<html lang="zh-TW">
		<head>
			<meta class="a b" name="description" itemprop="description" content="1"/>
			<meta class="a" name="description" itemprop="description" content="2"/>
			<meta class="a b" name="description" content="3"/>
			<meta class="a b" itemprop="description" content="4"/>
		</head>
		<body>
			<meta class="a b" name="description" itemprop="description" content="5"/>
		</body>
	</html>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmldoc))
	require.NoError(t, err)

	type testCase struct {
		Builder      *parser.SelectorBuilder
		NChild       int
		SelectedNode []string
	}

	tcs := []testCase{
		{
			Builder: parser.NewBuilder().
				Append(
					"head",
					nil,
					nil,
				).
				Append(
					"meta",
					[]string{"a", "b"},
					map[string]string{
						"name":     "description",
						"itemprop": "description",
					},
				),
			NChild:       1,
			SelectedNode: []string{"1"},
		},
		{
			Builder: parser.NewBuilder().
				Append(
					"head",
					nil,
					nil,
				).
				Append(
					"meta",
					[]string{"a"},
					map[string]string{
						"name":     "description",
						"itemprop": "description",
					},
				),
			NChild:       2,
			SelectedNode: []string{"1", "2"},
		},
		{
			Builder: parser.NewBuilder().
				Append(
					"head",
					nil,
					nil,
				).
				Append(
					"meta",
					[]string{"a", "b"},
					map[string]string{
						"itemprop": "description",
					},
				),
			NChild:       2,
			SelectedNode: []string{"1", "4"},
		},
		{
			Builder: parser.NewBuilder().
				Append(
					"head",
					nil,
					nil,
				).
				Append(
					"meta",
					[]string{"a", "b"},
					map[string]string{
						"name": "description",
					},
				),
			NChild:       2,
			SelectedNode: []string{"1", "3"},
		},
		{
			Builder: parser.NewBuilder().
				Append(
					"meta",
					[]string{"a", "b"},
					map[string]string{
						"name":     "description",
						"itemprop": "description",
					},
				),
			NChild:       2,
			SelectedNode: []string{"1", "5"},
		},
		{
			Builder: parser.NewBuilder().
				Append(
					"meta",
					[]string{"c"},
					map[string]string{
						"name":     "description",
						"itemprop": "description",
					},
				),
			NChild:       0,
			SelectedNode: []string{},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("Case %d", i+1),
			func(t *testing.T) {
				node := []string{}
				doc.Find(tc.Builder.Build()).Each(func(i int, s *goquery.Selection) {
					node = append(node, s.AttrOr("content", "X"))
				})
				require.ElementsMatch(t, tc.SelectedNode, node)
				require.Equal(t, tc.NChild, doc.Find(tc.Builder.Build()).Length())
			},
		)
	}

}

func TestFind(t *testing.T) {
	data := `[
    {
        "@context": "https://schema.org",
        "@type": "NewsArticle",
        "articleSection": "technology",
        "mainEntityOfPage": {
            "@type": "WebPage",
            "@id": "https://example.com/news/001.html"
        },
        "headline": "headline",
        "image": {
            "@type": "ImageObject",
            "url": "https://example.com/static/001.jpg",
            "height": 768,
            "width": 1024
        },
        "author": {
            "@type": "Person",
            "name": "John"
        },
        "url": "https://example.com/news/001.html",
        "thumbnailUrl": "https://example.com/static/001.jpg",
        "description": "description",
        "keywords": "kw1,kw2,kw3,kw4",
        "datePublished": "2023-08-25T19:04:00+08:00",
        "dateModified": "2023-08-25T19:04:00+08:00",
        "publisher": {
            "@type": "Organization",
            "name": "Organization Name",
            "logo": {
                "@type": "ImageObject",
                "url": "https://example.com/static/logo.jpg",
                "width": 260,
                "height": 60
            }
        }
    },
    {
        "@context": "http://schema.org",
        "@type": "Article",
        "url": "https://example.com/news/001.html",
        "thumbnailUrl": "https://example.com/static/001.jpg",
        "mainEntityOfPage": "https://example.com/news/001.html",
        "headline": "headline",
        "datePublished": "2023-08-25T19:04:00+08:00",
        "dateModified": "2023-08-25T19:04:00+08:00",
        "keywords": "kw1,kw2,kw3,kw4",
        "image": {
            "@type": "ImageObject",
            "url": "https://example.com/static/001.jpg",
            "height": 768,
            "width": 1024
        },
        "author": {
            "@type": "Person",
            "name": "John"
        },
        "publisher": {
            "@type": "Organization",
            "name": "Organization Name",
            "logo": {
                "@type": "ImageObject",
                "url": "https://example.com/static/logo.jpg",
                "width": 260,
                "height": 60
            }
        }
        "description": "description"
    }
]`
	p := parser.JsonLDParser{}
	newsa, err := p.FindTargetJsonLD([]byte(data), "NewsArticle")
	require.NoError(t, err)
	require.NotNil(t, newsa)

	require.True(t, json.Valid(newsa))

	var jld parser.JsonLD
	err = json.Unmarshal(newsa, &jld)
	require.NoError(t, err)
}
