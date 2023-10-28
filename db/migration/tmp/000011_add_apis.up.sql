INSERT INTO apis (
    id, name, type, image, icon, document_url
) VALUES 
    (1, 'NEWSDATA.IO' , 'source'        , 'logo_NEWSDATA.IO.png'          , 'favicon_NEWSDATA.IO.png', 'https://newsdata.io/documentation/'),
    (2, 'GNews'       , 'source'        , 'logo_GNews.png'                , 'favicon_GNews.ico'      , 'https://gnews.io/docs/v4'),
    (3, 'NEWS API'    , 'source'        , 'logo_NEWS_API.png'             , 'favicon_NEWS_API.ico'   , 'https://newsapi.org/docs/'),
    (4, 'OpenAI'      , 'language_model', 'logo_ChatGPT.svg'              , 'favicon_ChatGPT.ico'    , 'https://openai.com/blog/introducing-chatgpt-and-whisper-apis'),
    (5, 'Google API'  , 'source'        , 'logo_Google_Custom_Search.png' , 'favicon-Google.png'     , 'https://developers.google.com/custom-search/v1/reference/rest/v1/cse/list');

ALTER SEQUENCE apis_id_seq RESTART WITH 6;