package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	mock_model "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model/mockdb"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model/testtool"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	val "github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateAPIKeyService(t *testing.T) {
	type testCase struct {
		Name      string
		SetupFunc func(ctl *gomock.Controller) (*service.APIKeyCreateRequest, *mock_model.MockStore, int32)
		TestFunc  func(req *service.APIKeyCreateRequest, store *mock_model.MockStore, id int32)
	}

	Validate := val.New()
	validator.RegisterUUID(Validate)
	validator.RegisterValidator(
		Validate,
		validator.EnmusApiType,
	)

	tcs := []testCase{
		{
			Name: "Create API OK",
			SetupFunc: func(ctl *gomock.Controller) (*service.APIKeyCreateRequest, *mock_model.MockStore, int32) {
				apikey, err := testtool.GenRdmAPIKey(uuid.Nil, 0)
				require.NoError(t, err)

				req := &service.APIKeyCreateRequest{
					Owner: apikey.Owner,
					ApiID: apikey.ApiID,
					Key:   apikey.Key,
				}
				require.NotEmpty(t, req.RequestName())
				params, err := req.ToParams()
				require.NoError(t, err)

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					CreateAPIKey(gomock.Any(), gomock.Eq(params)).
					Times(1).
					Return(apikey.ID, nil)
				return req, store, apikey.ID
			},
			TestFunc: func(req *service.APIKeyCreateRequest, store *mock_model.MockStore, id int32) {
				srvc := service.NewService(store, Validate)
				n, err := srvc.APIKey().Create(context.Background(), req)
				require.NoError(t, err)
				require.Equal(t, id, n)
			},
		},
		{
			Name: "Missing Key",
			SetupFunc: func(ctl *gomock.Controller) (*service.APIKeyCreateRequest, *mock_model.MockStore, int32) {
				apikey, err := testtool.GenRdmAPIKey(uuid.Nil, 0)
				require.NoError(t, err)

				req := &service.APIKeyCreateRequest{
					Owner: apikey.Owner,
					ApiID: apikey.ApiID,
					Key:   "",
				}
				require.NotEmpty(t, req.RequestName())
				params, err := req.ToParams()
				require.NoError(t, err)

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					CreateAPI(gomock.Any(), gomock.Eq(params)).
					Times(0)

				return req, store, apikey.ID
			},
			TestFunc: func(req *service.APIKeyCreateRequest, store *mock_model.MockStore, id int32) {
				srvc := service.NewService(store, Validate)
				_, err := srvc.APIKey().Create(context.Background(), req)
				var valErr val.ValidationErrors
				ok := errors.As(err, &valErr)
				require.True(t, ok)
				require.Error(t, err)
				require.ErrorContains(t, err, "failed on the 'required' tag")
			},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("Case %d-%s", i+1, tc.Name),
			func(t *testing.T) {
				for j := 0; j < 10; j++ {
					func() {
						ctl := gomock.NewController(t)
						defer ctl.Finish()
						srvc, store, id := tc.SetupFunc(ctl)
						tc.TestFunc(srvc, store, id)
					}()
				}
			},
		)
	}
}

func TestListAPIKeyService(t *testing.T) {
	nUser := 5
	nAPI := 3

	users := make(map[uuid.UUID]*model.User, nUser)
	for i := 0; i < nUser; i++ {
		user, _ := testtool.GenRdmUser()
		users[user.ID] = user
	}

	apis := make(map[int16]*model.Api, nAPI)
	for i := 0; i < nAPI; i++ {
		api, _ := testtool.GenRdmAPI(-1)
		apis[api.ID] = api
	}

	type key struct {
		owner uuid.UUID
		api   int16
	}

	apikeys := make(map[key]*model.Apikey)
	for _, user := range users {
		for _, api := range apis {
			ak, _ := testtool.GenRdmAPIKey(user.ID, api.ID)
			apikeys[key{user.ID, api.ID}] = ak
		}
	}

	Validate := val.New()
	validator.RegisterUUID(Validate)
	for _, user := range users {
		func(user *model.User) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			store := mock_model.NewMockStore(ctl)
			store.
				EXPECT().
				ListAPIKey(gomock.Any(), gomock.Eq(user.ID)).
				Times(1).
				DoAndReturn(func(ctx context.Context, owner uuid.UUID) ([]*model.ListAPIKeyRow, error) {
					aks := []*model.ListAPIKeyRow{}
					for api_id, api := range apis {
						if apikey, ok := apikeys[key{owner, api_id}]; ok {
							aks = append(aks, &model.ListAPIKeyRow{
								ApiKeyID: apikey.ID,
								Owner:    owner,
								Key:      apikey.Key,
								ApiID:    api_id,
								Name:     api.Name,
								Type:     api.Type,
							})
						}
					}
					if len(aks) == 0 {
						return nil, pgx.ErrNoRows
					}
					return aks, nil
				})
			srvc := service.NewService(store, Validate)
			aks, err := srvc.APIKey().List(context.Background(), user.ID)
			require.NoError(t, err)
			for _, ak := range aks {
				require.Equal(t, user.ID, ak.Owner)
				require.Equal(t, apis[ak.ApiID].Name, ak.Name)
				require.Equal(t, apis[ak.ApiID].Type, ak.Type)
				require.Equal(t, apikeys[key{ak.Owner, ak.ApiID}].Key, ak.Key)
			}
		}(user)
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	invalidExistUserID, _ := uuid.NewDCEPerson() // uuid2
	store := mock_model.NewMockStore(ctl)
	store.
		EXPECT().
		ListAPIKey(gomock.Any(), gomock.Any()).
		Times(0)
	srvc := service.NewService(store, Validate)
	_, err := srvc.APIKey().List(context.Background(), invalidExistUserID)
	require.Error(t, err)

	var valErr val.ValidationErrors
	ok := errors.As(err, &valErr)
	require.True(t, ok)
	require.ErrorContains(t, err, "failed on the 'uuid4' tag")
}

func TestDeleteAPIKeyService(t *testing.T) {
	type ErrType int8
	const (
		Success ErrType = iota
		ValidationError
		NoRowError
	)

	type testCase struct {
		Name    string
		Owner   uuid.UUID
		ApiID   int16
		ErrType ErrType
		ErrTag  string
	}

	uuidv4 := uuid.New()
	uuidv2, err := uuid.NewDCEGroup()
	require.NoError(t, err)
	tcs := []testCase{
		{"OK", uuidv4, 1, Success, ""},
		{"Invalid Owner", uuidv2, 1, ValidationError, "uuid4"},
		{"Missing Owner", uuid.Nil, 1, ValidationError, "not_uuid_nil"},
		{"Invalid ApiID", uuidv4, -1, ValidationError, "min"},
		{"Missing ApiID", uuidv4, 0, ValidationError, "required"},
		{"No row", uuidv4, 1, NoRowError, ""},
	}

	Validate := val.New()
	validator.RegisterUUID(Validate)
	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.Name,
			func(t *testing.T) {
				ctl := gomock.NewController(t)
				defer ctl.Finish()
				store := mock_model.NewMockStore(ctl)

				switch tc.ErrType {
				case ValidationError:
					store.
						EXPECT().
						DeleteAPIKey(gomock.Any(), gomock.Any()).
						Times(0)
				case Success:
					store.
						EXPECT().
						DeleteAPIKey(
							gomock.Any(),
							gomock.Eq(&model.DeleteAPIKeyParams{
								Owner: tc.Owner,
								ApiID: tc.ApiID,
							})).Times(1).
						Return(int64(1), nil)
				case NoRowError:
					store.
						EXPECT().
						DeleteAPIKey(
							gomock.Any(),
							gomock.Eq(&model.DeleteAPIKeyParams{
								Owner: tc.Owner,
								ApiID: tc.ApiID,
							})).Times(1).
						Return(int64(0), pgx.ErrNoRows)
				}

				srvc := service.NewService(store, Validate)
				n, err := srvc.APIKey().Delete(
					context.Background(), &service.APIKeyDeleteRequest{
						Owner: tc.Owner,
						ApiID: tc.ApiID,
					})

				switch tc.ErrType {
				case ValidationError:
					require.Error(t, err)
					var valErr val.ValidationErrors
					ok := errors.As(err, &valErr)
					require.True(t, ok)
					require.ErrorContains(t, err, fmt.Sprintf("failed on the '%s' tag", tc.ErrTag))
				case Success:
					require.NoError(t, err)
					require.Equal(t, int64(1), n)
				case NoRowError:
					require.ErrorIs(t, err, pgx.ErrNoRows)
				}
			},
		)
	}
}
