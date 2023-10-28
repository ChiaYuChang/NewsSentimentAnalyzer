INSERT INTO endpoints (
    "name", "api_id", "template_name"
) VALUES
    ('Latest News'   , 1, 'NEWSDATA.IO-latest_news.gotmpl'),
    ('News Archive'  , 1, 'NEWSDATA.IO-news_archive.gotmpl'),
    ('News Sources'  , 1, 'NEWSDATA.IO-news_sources.gotmpl'),
    ('Search'        , 2, 'GNews-search.gotmpl'),
    ('Top Headlines' , 2, 'GNews-top_headlines.gotmpl'),
    ('Everything'    , 3, 'NewsAPI-everything.gotmpl'),
    ('Top Headlines' , 3, 'NewsAPI-top_headlines.gotmpl'),
    ('Sources'       , 3, 'NewsAPI-sources.gotmpl'),
    ('Custom Search' , 5, 'GoogleCSE.gotmpl');