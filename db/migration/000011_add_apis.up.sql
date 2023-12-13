INSERT INTO apis (
    id, name, type, image, icon, document_url, created_at, updated_at
) VALUES 
    (1, 'NEWSDATA.IO', 'source', 'logo_NEWSDATA.IO.png', 'favicon_NEWSDATA.IO.png', 'https://newsdata.io/documentation/', '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
    (2, 'GNews', 'source', 'logo_GNews.png', 'favicon_GNews.ico', 'https://gnews.io/docs/v4', '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
    (3, 'NEWS API', 'source', 'logo_NEWS_API.png', 'favicon_NEWS_API.ico', 'https://newsapi.org/docs/', '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
    (4, 'Google API', 'source', 'logo_Google_Custom_Search.png', 'favicon-Google.png', 'https://developers.google.com/custom-search/v1/reference/rest/v1/cse/list', '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
    (5, 'OpenAI', 'language_model', 'logo_ChatGPT.svg', 'favicon_ChatGPT.ico', 'https://openai.com/blog/introducing-chatgpt-and-whisper-apis', '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
    (6, 'Cohere', 'language_model', 'logo_Cohere.png', 'favicon_Cohere.png', 'https://cohere.com/', '2020-01-01 00:00:00', '2020-01-01 00:00:00');

ALTER SEQUENCE apis_id_seq RESTART WITH 7;