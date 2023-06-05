package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model/testtool"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
	"github.com/jackc/pgx/v5"
)

type CreateUpdateThenDeletAPI struct {
	n int
}

func (scrpt CreateUpdateThenDeletAPI) Do(
	ctx context.Context, srvc service.Service, logger Logger) Job {
	var err error
	var n int64
	jid := ctx.Value("jid").(int)

	logger.Print("\t- Check Create API...")
	apiOri, _ := testtool.GenRdmAPI()
	apiCreateReq := &service.APICreateRequest{
		Name: apiOri.Name,
		Type: string(apiOri.Type),
	}

	apiOri.ID, err = srvc.API().Create(context.Background(), apiCreateReq)
	if err != nil {
		logger.Println("Failed")
		return Job{jid, err}
	}
	apiInDB, _ := srvc.API().Get(context.Background(), apiOri.ID)
	if err != nil {
		logger.Println("Failed")
		return Job{jid, err}
	}
	if apiOri.Name != apiInDB.Name || apiOri.Type != apiInDB.Type {
		logger.Println("Failed")
		return Job{jid, ErrorFieldNotMatch}
	}
	logger.Println("OK")

	logger.Print("\t- Update API Name...")
	apiNameUpdated := testtool.CloneAPI(apiOri)
	apiNameUpdated.Name = rg.Must[string](rg.AlphaNum.GenRdmString(rand.Intn(10) + 10))
	apiUpdateNameReq := &service.APIUpdateRequeset{
		Name: apiNameUpdated.Name,
		Type: string(apiNameUpdated.Type),
		ID:   apiNameUpdated.ID,
	}

	n, err = srvc.API().Update(context.Background(), apiUpdateNameReq)
	if n != 1 {
		return Job{jid, fmt.Errorf("%w: expected: %d, actual: %d",
			ErrorAffectMoreThanExpectedRows, 1, n)}
	}
	if err != nil {
		logger.Println("Failed")
		return Job{jid, err}
	}

	apiInDB, _ = srvc.API().Get(context.Background(), apiOri.ID)
	if apiNameUpdated.Name != apiInDB.Name || apiNameUpdated.Type != apiInDB.Type {
		logger.Println("Failed")
		return Job{jid, ErrorFieldNotMatch}
	}
	logger.Println("OK")

	logger.Print("\t- Update API Type...")
	apiTypeUpdated := testtool.CloneAPI(apiNameUpdated)
	apiTypeUpdated.Name = rg.Must[string](rg.AlphaNum.GenRdmString(rand.Intn(10) + 10))
	apiTypeUpdatedReq := &service.APIUpdateRequeset{
		Name: apiTypeUpdated.Name,
		Type: string(apiTypeUpdated.Type),
		ID:   apiTypeUpdated.ID,
	}

	n, err = srvc.API().Update(context.Background(), apiTypeUpdatedReq)
	if n != 1 {
		return Job{jid, fmt.Errorf("%w: expected: %d, actual: %d",
			ErrorAffectMoreThanExpectedRows, 1, n)}
	}
	if err != nil {
		logger.Println("Failed")
		return Job{jid, err}
	}

	apiInDB, _ = srvc.API().Get(context.Background(), apiOri.ID)
	if apiTypeUpdated.Name != apiInDB.Name || apiTypeUpdated.Type != apiInDB.Type {
		logger.Println("Failed")
		return Job{jid, ErrorFieldNotMatch}
	}
	logger.Println("OK")

	logger.Print("\t- Delete API...")
	n, err = srvc.API().Delete(context.Background(), apiOri.ID)
	if n != 1 {
		return Job{jid, fmt.Errorf("%w: expected: %d, actual: %d",
			ErrorAffectMoreThanExpectedRows, 1, n)}
	}
	if err != nil {
		logger.Println("Failed")
		return Job{jid, err}
	}

	_, err = srvc.API().Get(context.Background(), apiOri.ID)
	if !errors.Is(err, pgx.ErrNoRows) {
		return Job{jid, err}
	}
	logger.Println("OK")

	return Job{jid, nil}
}

func (scrpt CreateUpdateThenDeletAPI) Steps() []string {
	return []string{
		"Create API",
		"Update API Name",
		"Update API Type",
		"Delete API",
	}
}

func (scrpt CreateUpdateThenDeletAPI) Description() string {
	return strings.Join(scrpt.Steps(), " -> ")
}

func (scrpt CreateUpdateThenDeletAPI) N() int {
	return scrpt.n
}
