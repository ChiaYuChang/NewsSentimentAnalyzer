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

	s := NewSampler(API_KEY_CHAR_SET, API_KEY_CHAR_PROB)
	for _, u := range users {
		for _, a := range apis {
			if (u.Id == TEST_ADMIN_USER_UID || u.Id == TEST_USER_UID) || rand.Float64() < a.Probability {
				apikeys.Item = append(apikeys.Item, APIKeyItem{
					Owner: u.Id,
					APIId: a.Id,
					Key:   string(s.GetN(API_KEY_LENGTH[rand.Intn(len(API_KEY_LENGTH))])),
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
