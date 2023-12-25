package main

import (
	"math/rand"
	"sort"
)

type NewsJobList []NewsJobItem

func (l NewsJobList) N() int {
	return len(l)
}

func (l NewsJobList) Item() []NewsJobItem {
	return l
}

type NewsJobItem struct {
	JobId  int64 `json:"job_id"`
	NewsId int64 `json:"news_id"`
}

func NewRandomNewsJob(n int, jids, nids []int64) NewsJobList {
	list := make(NewsJobList, 0, n)
	for i := 0; i < n; i++ {
		list = append(list, NewsJobItem{
			JobId:  jids[rand.Intn(len(jids))],
			NewsId: nids[rand.Intn(len(nids))],
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].JobId < list[j].JobId
	})

	return list
}
