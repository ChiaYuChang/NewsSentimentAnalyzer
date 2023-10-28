package main

import (
	"math/rand"
	"time"

	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
	"github.com/google/uuid"
)

type APIKey struct {
	Item []APIKeyItem
	N    int
}

type APIKeyItem struct {
	Id        int       `json:"id"                                    validate:"required"`
	Owner     uuid.UUID `json:"owner"                                 validate:"required"`
	APIId     int       `json:"api_id"                                validate:"required"`
	Key       string    `json:"key"        mod:"default=xxxx-xxx-xx"  validate:"required"`
	CreatedAt time.Time `json:"created_at" mod:"default=2000-01-01T00:00:00+00:00"`
	UpdatedAt time.Time `json:"updated_at" mod:"default=2000-01-01T00:00:00+00:00"`
}

func NewApiKeys(apis []APIItem, users []UserItem) APIKey {
	apikeys := APIKey{}

	kl := []int{16, 32, 52}
	s := NewSampler([]rune{'0', 'x', '-'}, []float64{0.6, 0.35, 0.05})
	for _, u := range users {
		for _, a := range apis {
			if rand.Float64() < a.Probability {
				apikeys.Item = append(apikeys.Item, APIKeyItem{
					Owner: u.Id,
					APIId: a.Id,
					Key:   string(s.GetN(kl[rand.Intn(len(kl))])),
				})
			}
		}
	}

	rand.Shuffle(len(apikeys.Item), func(i, j int) {
		apikeys.Item[i], apikeys.Item[j] = apikeys.Item[j], apikeys.Item[i]
	})

	rts := rg.GenRdnTimes(len(apikeys.Item), TIME_MIN, TIME_MAX)
	for i := range apikeys.Item {
		apikeys.Item[i].Id = i + 1
		apikeys.Item[i].CreatedAt = rts[i]
		apikeys.Item[i].UpdatedAt = rg.GenRdnTime(rts[i], TIME_MAX)
	}

	apikeys.N = len(apikeys.Item) + 1
	return apikeys
}
