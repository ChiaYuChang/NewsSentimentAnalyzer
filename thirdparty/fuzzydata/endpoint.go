package main

import "time"

type EndpointList []EndpointItem

func (e EndpointList) N() int {
	return len(e) + 1
}

type EndpointItem struct {
	Name         string    `json:"name"          mod:"trim"`
	APIName      string    `json:"api_name"      mod:"trim"`
	APIId        int       `json:"api_id"`
	TemplateName string    `json:"template_name" mod:"trim"`
	CreatedAt    time.Time `json:"created_at"    mod:"default=2000-01-01T00:00:00+00:00"`
	UpdatedAt    time.Time `json:"updated_at"    mod:"default=2000-01-01T00:00:00+00:00"`
}

var Endpoints = EndpointList{
	{
		Name:         "Latest News",
		APIId:        1,
		TemplateName: "NEWSDATA.IO-latest_news.gotmpl",
		CreatedAt:    TIME_MIN,
		UpdatedAt:    TIME_MIN,
	},
	// {
	// 	Name:         "News Archive",
	// 	APIId:        1,
	// 	TemplateName: "NEWSDATA.IO-news_archive.gotmpl",
	// 	CreatedAt:    TIME_MIN,
	// 	UpdatedAt:    TIME_MIN,
	// },
	// {
	// 	Name:         "News Sources",
	// 	APIId:        1,
	// 	TemplateName: "NEWSDATA.IO-news_sources.gotmpl",
	// 	CreatedAt:    TIME_MIN,
	// 	UpdatedAt:    TIME_MIN,
	// },
	{
		Name:         "Search",
		APIId:        2,
		TemplateName: "GNews-search.gotmpl",
		CreatedAt:    TIME_MIN,
		UpdatedAt:    TIME_MIN,
	},
	{
		Name:         "Top Headlines",
		APIId:        2,
		TemplateName: "GNews-top_headlines.gotmpl",
		CreatedAt:    TIME_MIN,
		UpdatedAt:    TIME_MIN,
	},
	{
		Name:         "Everything",
		APIId:        3,
		TemplateName: "NewsAPI-everything.gotmpl",
		CreatedAt:    TIME_MIN,
		UpdatedAt:    TIME_MIN,
	},
	{
		Name:         "Top Headlines",
		APIId:        3,
		TemplateName: "NewsAPI-top_headlines.gotmpl",
		CreatedAt:    TIME_MIN,
		UpdatedAt:    TIME_MIN,
	},
	// {
	// 	Name:         "Sources",
	// 	APIId:        3,
	// 	TemplateName: "NewsAPI-sources.gotmpl",
	// 	CreatedAt:    TIME_MIN,
	// 	UpdatedAt:    TIME_MIN,
	// },
	{
		Name:         "Custom Search",
		APIId:        4,
		TemplateName: "GoogleCSE.gotmpl",
		CreatedAt:    TIME_MIN,
		UpdatedAt:    TIME_MIN,
	},
}
