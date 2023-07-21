package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	mock_model "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model/mockdb"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model/testtool"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
	val "github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestListEndpointByOwener(t *testing.T) {
	type testCase struct {
		Name      string
		Owner     uuid.UUID
		SetupFunc func(owner uuid.UUID, ctl *gomock.Controller) service.Service
		CheckFunc func(owner uuid.UUID, srvc service.Service)
	}

	n := 5
	user, _ := testtool.GenRdmUser()
	apis := make([]*model.Api, n)
	for i := 0; i < n; i++ {
		apis[i], _ = testtool.GenRdmAPI(int16(i))
	}

	apikeys := make([]*model.Apikey, n)
	for i := 0; i < n; i++ {
		apikeys[i], _ = testtool.GenRdmAPIKey(user.ID, apis[i].ID)
	}

	eps := make([]*model.Endpoint, n)
	for i := 0; i < n; i++ {
		eps[i], _ = testtool.GenRdmAPIEndpoints(-1, apis[i].ID)
	}

	rows := make([]*model.ListEndpointByOwnerRow, n)
	for i := 0; i < n; i++ {
		rows[i] = &model.ListEndpointByOwnerRow{
			EndpointID:   eps[i].ID,
			EndpointName: eps[i].Name,
			ApiID:        eps[i].ApiID,
			TemplateName: eps[i].TemplateName,
			Key:          apikeys[i].Key,
			ApiName:      apis[i].Name,
			Type:         apis[i].Type,
			Icon:         apis[i].Icon,
			Image:        apis[i].Image,
			DocumentUrl:  apis[i].DocumentUrl,
		}
	}

	tcs := []testCase{
		{
			Name:  "OK",
			Owner: uuid.New(),
			SetupFunc: func(owner uuid.UUID, ctl *gomock.Controller) service.Service {
				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					ListEndpointByOwner(gomock.Any(), gomock.Eq(owner)).
					Times(1).
					Return(rows, nil)

				return service.NewService(store, validator.Validate)
			},
			CheckFunc: func(owner uuid.UUID, srvc service.Service) {
				eps, err := srvc.
					Endpoint().
					ListEndpointByOwner(
						context.Background(),
						owner)
				require.NoError(t, err)
				require.Equal(t, eps, rows)
			},
		},
		{
			Name:  "Error Owner ID",
			Owner: uuid.Nil,
			SetupFunc: func(owner uuid.UUID, ctl *gomock.Controller) service.Service {
				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					ListEndpointByOwner(gomock.Any(), gomock.Eq(owner)).
					Times(0)
				return service.NewService(store, validator.Validate)
			},
			CheckFunc: func(owner uuid.UUID, srvc service.Service) {
				eps, err := srvc.
					Endpoint().
					ListEndpointByOwner(
						context.Background(),
						owner)
				require.Error(t, err)
				var valErr val.ValidationErrors
				require.ErrorAs(t, err, &valErr)
				require.Nil(t, eps)
			},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		ctl := gomock.NewController(t)
		t.Run(
			tc.Name,
			func(t *testing.T) {
				t.Parallel()
				srvc := tc.SetupFunc(tc.Owner, ctl)
				tc.CheckFunc(tc.Owner, srvc)
			},
		)
	}
}

func TestCreateEndpoint(t *testing.T) {
	ep, _ := testtool.GenRdmAPIEndpoints(-1, -1)
	ep.Name = rg.Must[string](rg.Alphabet.GenRdmString(32))
	ep.TemplateName = rg.Must[string](rg.Alphabet.GenRdmString(25)) + ".gotmpl"

	type testCase struct {
		Name      string
		Request   *service.EndpointCreateRequest
		SetupFunc func(req *service.EndpointCreateRequest, ctl *gomock.Controller) service.Service
		CheckFunc func(req *service.EndpointCreateRequest, srvc service.Service)
	}

	tcs := []testCase{
		{
			Name: "OK",
			Request: &service.EndpointCreateRequest{
				Name:         ep.Name,
				ApiID:        ep.ApiID,
				TemplateName: ep.TemplateName,
			},
			SetupFunc: func(req *service.EndpointCreateRequest, ctl *gomock.Controller) service.Service {
				param, _ := req.ToParams()
				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					CreateEndpoint(gomock.Any(), gomock.Eq(param)).
					Times(1).
					Return(ep.ID, nil)
				return service.NewService(store, validator.Validate)
			},
			CheckFunc: func(req *service.EndpointCreateRequest, srvc service.Service) {
				id, err := srvc.Endpoint().Create(context.Background(), req)
				require.NoError(t, err)
				require.Equal(t, ep.ID, id)
			},
		},
		{
			Name: "Endpoint name too long",
			Request: &service.EndpointCreateRequest{
				Name:         ep.Name + "a",
				ApiID:        ep.ApiID,
				TemplateName: ep.TemplateName,
			},
			SetupFunc: func(req *service.EndpointCreateRequest, ctl *gomock.Controller) service.Service {
				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					CreateEndpoint(gomock.Any(), gomock.Any()).
					Times(0)
				return service.NewService(store, validator.Validate)
			},
			CheckFunc: func(req *service.EndpointCreateRequest, srvc service.Service) {
				id, err := srvc.Endpoint().Create(context.Background(), req)
				require.Error(t, err)
				var valErr val.ValidationErrors
				require.ErrorAs(t, err, &valErr)
				require.Equal(t, int32(0), id)
				require.ErrorContains(t, err, "'Name' failed on the 'max' tag")
			},
		},
		{
			Name: "Template name too long",
			Request: &service.EndpointCreateRequest{
				Name:         ep.Name,
				ApiID:        ep.ApiID,
				TemplateName: ep.TemplateName + "a",
			},
			SetupFunc: func(req *service.EndpointCreateRequest, ctl *gomock.Controller) service.Service {
				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					CreateEndpoint(gomock.Any(), gomock.Any()).
					Times(0)
				return service.NewService(store, validator.Validate)
			},
			CheckFunc: func(req *service.EndpointCreateRequest, srvc service.Service) {
				id, err := srvc.Endpoint().Create(context.Background(), req)
				require.Error(t, err)
				var valErr val.ValidationErrors
				require.ErrorAs(t, err, &valErr)
				require.Equal(t, int32(0), id)
				require.ErrorContains(t, err, "'TemplateName' failed on the 'max' tag")
			},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		ctl := gomock.NewController(t)
		t.Run(
			tc.Name,
			func(t *testing.T) {
				t.Parallel()
				require.NotEmpty(t, tc.Request.RequestName())
				srvc := tc.SetupFunc(tc.Request, ctl)
				tc.CheckFunc(tc.Request, srvc)
			},
		)
	}
}

func TestListAllEndpoint(t *testing.T) {
	n := 30
	eps := make([]*model.Endpoint, n)
	for i := 1; i <= n; i++ {
		eps[i-1], _ = testtool.GenRdmAPIEndpoints(int32(i), -1)
	}

	tcs := []int{7, 10, 11}

	for i, limit := range tcs {

		t.Run(
			fmt.Sprintf("Case %d-limit=%d", i+1, limit),
			func(t *testing.T) {
				times := n / limit
				if n%limit != 0 {
					times++
				}
				times++
				ctl := gomock.NewController(t)
				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					ListAllEndpoint(gomock.Any(), gomock.Any()).
					Times(times).
					DoAndReturn(func(ctx context.Context, params *model.ListAllEndpointParams) ([]*model.ListAllEndpointRow, error) {
						if int(params.Next) >= len(eps) {
							return nil, pgx.ErrNoRows
						}

						rows := make([]*model.ListAllEndpointRow, 0, params.Limit)
						for _, ep := range eps {
							if ep.ID <= params.Next {
								continue
							}

							rows = append(rows, &model.ListAllEndpointRow{
								EndpointID:   ep.ID,
								EndpointName: ep.Name,
								ApiID:        ep.ApiID,
								TemplateName: ep.TemplateName,
							})

							if len(rows) >= int(params.Limit) {
								break
							}
						}
						return rows, nil
					})

				rowChan := make(chan *model.ListAllEndpointRow)
				go func(ch chan *model.ListAllEndpointRow) {
					srvc := service.NewService(store, validator.Validate)
					srvc.Endpoint().ListAll(context.Background(), int32(limit), ch)
				}(rowChan)

				i := int32(1)
				for row := range rowChan {
					require.Equal(t, row.EndpointID, i)
					i++
				}
				require.Equal(t, int32(n+1), i)
			},
		)
	}
}
