package gnews

type Response struct {
	TotalArticles int       `json:"totalArticles"`
	Articles      []Article `json:"articles"`
}

type Article struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Content     string        `json:"content"`
	URL         string        `json:"url"`
	Image       string        `json:"image"`
	PublishedAt string        `json:"publishedAt"`
	Source      ArticleSource `json:"source"`
}

type ArticleSource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
