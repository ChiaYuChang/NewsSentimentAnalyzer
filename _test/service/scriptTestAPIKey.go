package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model/testtool"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/jackc/pgx/v5"
)

type CreateGetThenDeletAPIKey struct {
	n     int
	nUser int
	nAPI  int
}

func (scrpt CreateGetThenDeletAPIKey) Do(
	ctx context.Context, srvc service.Service, logger Logger) Job {
	var err error
	var n int64
	jid := ctx.Value("jid").(int)

	logger.Print("\t- Create User and API...")
	users := make([]int32, scrpt.nUser)
	for i := 0; i < scrpt.nUser; i++ {
		user, _ := testtool.GenRdmUser()
		users[i], _ = srvc.User().Create(ctx, &service.UserCreateRequest{
			Password:  string(user.Password),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
			Email:     user.Email,
		})
	}

	apis := make([]int16, scrpt.nAPI)
	for i := 0; i < scrpt.nAPI; i++ {
		api, _ := testtool.GenRdmAPI()
		apis[i], _ = srvc.API().Create(ctx, &service.APICreateRequest{
			Name: api.Name,
			Type: string(api.Type),
		})
	}
	logger.Println("OK")

	j := 0
	for iUser := 0; iUser < scrpt.nUser; iUser++ {
		for iAPI := 0; iAPI < scrpt.nAPI; iAPI++ {
			logger.Printf("\t- Iter: %02d/%02d...\n", j, scrpt.nUser*scrpt.nAPI)
			logger.Print("\t  - Check Create APIKey...")
			apikeyOri, _ := testtool.GenRdmAPIKey(
				users[iUser], apis[iAPI],
			)
			j++

			apikeyCreateReq := &service.APIKeyCreateRequest{
				Owner: apikeyOri.Owner,
				ApiID: apikeyOri.ApiID,
				Key:   apikeyOri.Key,
			}

			apikeyOri.ID, err = srvc.APIKey().Create(
				context.Background(),
				apikeyCreateReq,
			)
			if err != nil {
				logger.Println("Failed")
				return Job{jid, err}
			}
			logger.Println("OK")

			logger.Print("\t  - Check Get APIKey...")
			apikeyInDB, err := srvc.APIKey().Get(
				context.Background(),
				&service.APIKeyGetRequest{
					Owner: apikeyOri.Owner,
					ApiID: apikeyOri.ApiID,
				})
			if err != nil {
				logger.Println("Failed")
				return Job{jid, err}
			}

			if apikeyOri.ID != apikeyInDB.ID ||
				apikeyOri.ApiID != apikeyInDB.ApiID ||
				apikeyOri.Owner != apikeyInDB.Owner ||
				apikeyOri.Key != apikeyInDB.Key {
				logger.Println("Failed")
				return Job{jid, ErrorFieldNotMatch}
			}
			logger.Println("OK")

			logger.Print("\t  - Check Delete APIKey...")
			n, err = srvc.APIKey().Delete(context.Background(), &service.APIKeyDeleteRequest{
				Owner: apikeyOri.Owner,
				ApiID: apikeyOri.ApiID,
			})
			if n != 1 {
				return Job{jid, fmt.Errorf("%w: expected: %d, actual: %d",
					ErrorAffectMoreThanExpectedRows, 1, n)}
			}
			if err != nil {
				logger.Println("Failed")
				return Job{jid, err}
			}

			_, err = srvc.APIKey().Get(context.Background(), &service.APIKeyGetRequest{
				Owner: apikeyOri.Owner,
				ApiID: apikeyOri.ApiID,
			})
			if !errors.Is(err, pgx.ErrNoRows) {
				return Job{jid, err}
			}
			logger.Println("OK")
		}
	}
	return Job{jid, nil}
}

func (scrpt CreateGetThenDeletAPIKey) Steps() []string {
	return []string{
		"Create User and API",
		"Create APIKey",
		"Get APIKey Type",
		"Delete APIKey",
	}
}

func (scrpt CreateGetThenDeletAPIKey) Description() string {
	return strings.Join(scrpt.Steps(), " -> ")
}

func (scrpt CreateGetThenDeletAPIKey) N() int {
	return scrpt.n
}

type UserDeleteCascase struct {
	n int
}

func (srcpt UserDeleteCascase) Do(
	ctx context.Context, srvc service.Service, logger Logger) Job {
	var err error
	var n int64
	jid := ctx.Value("jid").(int)

	logger.Print("\t- Create User and API...")
	user, _ := testtool.GenRdmUser()
	user.ID, _ = srvc.User().Create(ctx, &service.UserCreateRequest{
		Password:  string(user.Password),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
		Email:     user.Email,
	})

	api, _ := testtool.GenRdmAPI()
	api.ID, _ = srvc.API().Create(ctx, &service.APICreateRequest{
		Name: api.Name,
		Type: string(api.Type),
	})
	logger.Println("OK")

	logger.Print("\t- Check Create APIKey...")
	apikeyOri, _ := testtool.GenRdmAPIKey(user.ID, api.ID)

	apikeyCreateReq := &service.APIKeyCreateRequest{
		Owner: apikeyOri.Owner,
		ApiID: apikeyOri.ApiID,
		Key:   apikeyOri.Key,
	}

	apikeyOri.ID, err = srvc.APIKey().Create(
		context.Background(),
		apikeyCreateReq,
	)
	if err != nil {
		logger.Println("Failed")
		return Job{jid, err}
	}
	logger.Println("OK")

	logger.Print("\t- Delete user...")
	n, err = srvc.User().HardDelete(ctx, user.ID)
	if err != nil {
		logger.Println("Failed")
		return Job{jid, err}
	}
	if n != 1 {
		return Job{jid, logger.Errorf("%w: expected %d, actual %d",
			ErrorAffectMoreThanExpectedRows, 1, n)}
	}
	logger.Println("OK")

	logger.Print("\t- Get Apikey (No Row)...")
	_, err = srvc.APIKey().Get(context.Background(),
		&service.APIKeyGetRequest{
			Owner: apikeyOri.Owner,
			ApiID: apikeyOri.ApiID,
		})
	if !errors.Is(err, pgx.ErrNoRows) {
		return Job{jid, err}
	}
	logger.Println("OK")
	return Job{jid, nil}
}

func (scrpt UserDeleteCascase) Steps() []string {
	return []string{
		"Create User and API",
		"Create APIKey",
		"Delete APIKey",
		"Get APIKey",
	}
}

func (scrpt UserDeleteCascase) Description() string {
	return strings.Join(scrpt.Steps(), " -> ")
}

func (scrpt UserDeleteCascase) N() int {
	return scrpt.n
}
