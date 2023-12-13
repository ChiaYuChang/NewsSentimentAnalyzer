package service

import (
	"context"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/google/uuid"
)

func (srvc txService) DoCacheToStoreTx(ctx context.Context, user uuid.UUID,
	ulid string, srcId int16, srcQuery string, llmId int16, llmQuery string,
	cnReqChan <-chan *NewsCreateRequest) (*model.CacheToStoreTXResult, error) {

	cnpChan := make(chan *model.CreateNewsParams, 10)
	go func() {
		defer close(cnpChan)
		for req := range cnReqChan {
			params, _ := req.ToParams()
			cnpChan <- params
		}
	}()

	return srvc.store.DoCacheToStoreTx(ctx, &model.CacheToStoreTXParams{
		CreateJobParams: &model.CreateJobParams{
			Owner:    user,
			Ulid:     ulid,
			SrcApiID: srcId,
			SrcQuery: srcQuery,
			LlmApiID: llmId,
			LlmQuery: []byte(llmQuery),
		},
		CreateNewsParams: cnpChan,
	})
}
