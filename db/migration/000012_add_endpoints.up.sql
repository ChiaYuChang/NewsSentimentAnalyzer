INSERT INTO endpoints (
    id, name, api_id, template_name, created_at, updated_at
) VALUES
    (1, 'Latest News', '1', 'NEWSDATA.IO-latest_news.gotmpl', '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
    (2, 'Search', '2', 'GNews-search.gotmpl', '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
    (3, 'Top Headlines', '2', 'GNews-top_headlines.gotmpl', '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
    (4, 'Everything', '3', 'NewsAPI-everything.gotmpl', '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
    (5, 'Top Headlines', '3', 'NewsAPI-top_headlines.gotmpl', '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
    (6, 'Sources', '3', 'NewsAPI-sources.gotmpl', '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
    (7, 'Custom Search', '4', 'GoogleCSE.gotmpl', '2020-01-01 00:00:00', '2020-01-01 00:00:00');

ALTER SEQUENCE endpoints_id_seq RESTART WITH 8;
