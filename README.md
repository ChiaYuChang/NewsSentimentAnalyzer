# 具圖形化的台灣媒體分析工具

## 1. 台灣媒體問題

- 媒體為了吸引點閱率會使用聳動但是意義不完整的標題
- 媒體對於新聞的選擇會有偏頗
- 報導的內容帶有特定立場

## 2. [Ground News](https://ground.news/)

- 將新聞媒體分為左中右派
- 提供每一則新聞被各個媒體報導的比率
- 總結左中右派觀點並相互比較
- 將媒體擁有者分為八類 (傳媒、獨立媒體、政府...)
- 提供客戶自身盲點的新聞

## 3. 目標

### 3-1. 建立圖形化界面取得文章

#### 3-1-1. 文章蒐集 APIs

- [NewsAPI](https://newsapi.org/)
- [GNews](https://gnews.io/)
- [NEWSDATA.IO](https://newsdata.io/)
- [Google Custom Search JSON API](https://developers.google.com/custom-search/v1/introduction)

#### 3-1-2. 目標網域

##### 國內媒體

- [x] ETtoday 新聞雲 <www.ettoday.net>
- [x] TVBS 新聞網 (<news.tvbs.com.tw>)
- [x] 三立新聞網 (<www.setn.com>)
- [x] 上報 (<www.upmedia.mg>)
- [x] 中央社 (<www.cna.com.tw>)
- [x] 中時新聞網 (<www.chinatimes.com>)
- [x] 公視新聞網 (<news.pts.org.tw>)
- [x] 自由時報 (<news.ltn.com.tw>)
- [x] 聯合新聞網 (<udn.com>)
- [x] 轉角國際 (<global.udn.com>)
- [ ] TechNews 科技新報 (<technews.tw>)

##### 外媒

- [x] 法廣 (RFI) 台灣 (<www.rfi.fr/tw>)
- [x] 紐約時報 (The New York Times) 中文網 (<cn.nytimes.com/zh-hant/>)
- [x] BBC News 中文 <www.bbc.com/zhongwen/trad>
- [ ] 德國之聲 中文 (<www.dw.com>)
- [ ] NHK World News 中文 <www3.nhk.or.jp/nhkworld/zt/>
- [ ] 新華社 (<www.xinhuanet.com>)

### 3-2. 文章分析

- 以 OpenAI 的 [embeddings API endpoint](https://platform.openai.com/docs/guides/embeddings/what-are-embeddings) 取得各媒體針對同一主題的新聞文章的 Embedding。
  - 觀察各媒體文章標題 (title) 是否有差異
  - 觀察各媒體文章內容 (content) 是否有差異
- 以大型語言模型 (large language model，LLM) 分析各媒體對特定主題的情緒 (Sentiment Analysis)
  - 將文章分為 5 類
    - 1 (Very Negative)
    - 2 (Negative)
    - 3 (Neutral)
    - 4 (Positive)
    - 5 (Very Positive)
  - 可能的 LLM
    - [ChatGPT](https://chat.openai.com/)
      - 可以直接輸入中文進行分析
      - [Chat completions API](https://platform.openai.com/docs/api-reference/chat)
      - [Completions API](https://platform.openai.com/docs/guides/gpt/completions-api)
    - [Claude](https://claude.ai/)
      - 可以直接輸入中文進行分析
    - [PrivateGPT](https://github.com/imartinez/privateGPT)
    - [LlamaGPT](https://github.com/getumbrel/llama-gpt)
    - [Taiwan-LLaMa](https://github.com/MiuLab/Taiwan-LLaMa)
      - 模型訓練時加入台灣資料集
