package service_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	mrand "math/rand"
	"sort"
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	mock_model "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model/mockdb"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model/testtool"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
	val "github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateAPIService(t *testing.T) {
	type testCase struct {
		Name        string
		Init        func()
		SetupFunc   func(ctl *gomock.Controller) (*service.APICreateRequest, *mock_model.MockStore, int16)
		TestFunc    func(req *service.APICreateRequest, store *mock_model.MockStore, id int16)
		CleanUpFunc func()
	}

	var Validate *val.Validate
	tcs := []testCase{
		{
			Name: "Create API OK",
			Init: func() {
				Validate = val.New()
				validator.RegisterValidator(
					Validate,
					validator.EnmusApiType,
				)
			},
			SetupFunc: func(ctl *gomock.Controller) (*service.APICreateRequest, *mock_model.MockStore, int16) {
				api, err := testtool.GenRdmAPI(-1)
				require.NoError(t, err)

				req := &service.APICreateRequest{
					Name: api.Name,
					Type: string(api.Type),
				}
				require.NotEmpty(t, req.RequestName())
				params, err := req.ToParams()
				require.NoError(t, err)

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					CreateAPI(gomock.Any(), gomock.Eq(params)).
					Times(1).
					Return(api.ID, nil)
				return req, store, api.ID
			},
			TestFunc: func(req *service.APICreateRequest, store *mock_model.MockStore, id int16) {
				srvc := service.NewService(store, Validate)
				n, err := srvc.API().Create(context.Background(), req)
				require.NoError(t, err)
				require.Equal(t, id, n)
			},
			CleanUpFunc: func() {
				Validate = nil
			},
		},
		{
			Name: "Bad role",
			Init: func() {
				Validate = val.New()
				validator.RegisterValidator(
					Validate,
					validator.EnmusApiType,
				)
			},
			SetupFunc: func(ctl *gomock.Controller) (*service.APICreateRequest, *mock_model.MockStore, int16) {
				api, err := testtool.GenRdmAPI(-1)
				require.NoError(t, err)

				req := &service.APICreateRequest{
					Name: api.Name,
					Type: rg.Must[string](rg.AlphaNum.GenRdmString(10)),
				}
				require.NotEmpty(t, req.RequestName())
				params, err := req.ToParams()
				require.NoError(t, err)

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					CreateAPI(gomock.Any(), gomock.Eq(params)).
					Times(0)

				return req, store, api.ID
			},
			TestFunc: func(req *service.APICreateRequest, store *mock_model.MockStore, id int16) {
				srvc := service.NewService(store, Validate)
				_, err := srvc.API().Create(context.Background(), req)
				var valErr val.ValidationErrors
				ok := errors.As(err, &valErr)
				require.True(t, ok)
				require.Error(t, err)
				require.ErrorContains(t, err, "failed on the 'api_type' tag")
			},
			CleanUpFunc: func() {
				Validate = nil
			},
		},
		{
			Name: "Bad API Name",
			Init: func() {
				Validate = val.New()
				validator.RegisterValidator(
					Validate,
					validator.EnmusApiType,
				)
			},
			SetupFunc: func(ctl *gomock.Controller) (*service.APICreateRequest, *mock_model.MockStore, int16) {
				api, err := testtool.GenRdmAPI(-1)
				require.NoError(t, err)

				l := mrand.Int63n(3)
				if mrand.Float64() > 0.5 {
					l += mrand.Int63n(10) + 20
				}
				req := &service.APICreateRequest{
					Name: rg.Must[string](rg.AlphaNum.GenRdmString(int(l))),
					Type: string(api.Type),
				}
				require.NotEmpty(t, req.RequestName())
				params, err := req.ToParams()
				require.NoError(t, err)

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					CreateAPI(gomock.Any(), gomock.Eq(params)).
					Times(0)

				return req, store, api.ID
			},
			TestFunc: func(req *service.APICreateRequest, store *mock_model.MockStore, id int16) {
				srvc := service.NewService(store, Validate)
				_, err := srvc.API().Create(context.Background(), req)
				var valErr val.ValidationErrors
				ok := errors.As(err, &valErr)
				require.True(t, ok)
				require.Error(t, err)
				if len(req.Name) < 1 {
					require.ErrorContains(t, err, "failed on the 'required' tag")
				} else if len(req.Name) < 3 {
					require.ErrorContains(t, err, "failed on the 'min' tag")
				} else {
					require.ErrorContains(t, err, "failed on the 'max' tag")
				}
			},
			CleanUpFunc: func() {
				Validate = nil
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
						if tc.Init != nil {
							tc.Init()
						}

						if tc.CleanUpFunc != nil {
							defer tc.CleanUpFunc()
						}

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

func TestListAPIService(t *testing.T) {
	tcs := []int32{0, 1, 10, 50, 100, 101}
	min := int32(1)
	max := int32(100)

	Validate := val.New()
	for _, limit := range tcs {
		ctl := gomock.NewController(t)
		store := mock_model.NewMockStore(ctl)

		var apis []*model.ListAPIRow
		if limit < min || limit > max {
			store.
				EXPECT().
				ListAPI(gomock.Any(), gomock.Eq(limit)).
				Times(0)
		} else {
			apis = make([]*model.ListAPIRow, limit)
			for i := range apis {
				tmp, _ := testtool.GenRdmAPI(-1)
				apis[i] = &model.ListAPIRow{
					ID:   tmp.ID,
					Name: tmp.Name,
					Type: model.ApiType(tmp.Name),
				}
			}
			sort.Slice(apis, func(i, j int) bool {
				if apis[i].Type != apis[j].Type {
					if apis[i].Type < apis[j].Type {
						return true
					} else {
						return false
					}
				}
				if apis[i].Name != apis[i].Name {
					if apis[i].Name < apis[j].Name {
						return true
					} else {
						return false
					}
				}
				return apis[i].ID < apis[j].ID
			})

			store.
				EXPECT().
				ListAPI(gomock.Any(), gomock.Eq(limit)).
				Times(1).
				Return(apis, nil)
		}

		srvc := service.NewService(store, Validate)
		if limit < min || limit > max {
			rows, err := srvc.API().List(context.Background(), int32(limit))
			var valErr val.ValidationErrors
			ok := errors.As(err, &valErr)
			require.True(t, ok)
			require.Error(t, err)
			require.Nil(t, rows)
		} else {
			rows, err := srvc.API().List(context.Background(), int32(limit))
			require.NoError(t, err)
			require.NotEmpty(t, rows)
		}
	}
}

// TODO func TestUpdateAPIService(t *testing.T)

func TestDeleteAPIService(t *testing.T) {
	type testCase struct {
		Name      string
		SetupFunc func(ctl *gomock.Controller) (int16, *mock_model.MockStore, int64)
		TestFunc  func(id int16, store *mock_model.MockStore, nAffectedRow int64)
	}

	Validate := val.New()
	tcs := []testCase{
		{
			Name: "Delete API OK",
			SetupFunc: func(ctl *gomock.Controller) (int16, *mock_model.MockStore, int64) {
				id := int16(rand.Int31n(10_000) + 1)

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					DeleteAPI(gomock.Any(), gomock.Eq(id)).
					Times(1).
					Return(int64(1), nil)

				return id, store, 1
			},
			TestFunc: func(id int16, store *mock_model.MockStore, nAffectedRow int64) {
				srvc := service.NewService(store, Validate)
				n, err := srvc.API().Delete(context.Background(), id)
				require.NoError(t, err)
				require.Equal(t, nAffectedRow, n)
			},
		},
		{
			Name: "ID is invalid",
			SetupFunc: func(ctl *gomock.Controller) (int16, *mock_model.MockStore, int64) {
				id := -int16(rand.Int31n(10_000))

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(0)

				return id, store, 0
			},
			TestFunc: func(id int16, store *mock_model.MockStore, nAffectedRow int64) {
				srvc := service.NewService(store, Validate)
				n, err := srvc.API().Delete(context.Background(), id)
				var valErr val.ValidationErrors
				ok := errors.As(err, &valErr)
				require.True(t, ok)
				require.Error(t, err)
				require.Equal(t, int64(0), n)
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

func TestCleanUpAPIService(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	nAffectedRow := rand.Int63n(1_000) + 1
	store := mock_model.NewMockStore(ctl)
	store.
		EXPECT().
		CleanUpAPIs(gomock.Any()).
		Times(1).
		Return(nAffectedRow, nil)

	srvc := service.NewService(store, nil)
	n, err := srvc.API().CleanUp(context.Background())
	require.NoError(t, err)
	require.Equal(t, nAffectedRow, n)
}
