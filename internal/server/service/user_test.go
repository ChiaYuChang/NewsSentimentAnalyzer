package service_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	mock_model "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model/mockdb"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model/testtool"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
	val "github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUserService(t *testing.T) {
	type testCase struct {
		Name        string
		Init        func()
		SetupFunc   func(ctl *gomock.Controller) (*service.UserCreateRequest, *mock_model.MockStore, int32)
		TestFunc    func(req *service.UserCreateRequest, store *mock_model.MockStore, id int32)
		CleanUpFunc func()
	}

	var Validate *val.Validate
	tcs := []testCase{
		{
			Name: "Create User OK",
			Init: func() {
				Validate = val.New()
				validator.RegisterValidator(
					Validate,
					validator.EnmusRole,
					validator.NewDefaultPasswordValidator(),
				)
			},
			SetupFunc: func(ctl *gomock.Controller) (*service.UserCreateRequest, *mock_model.MockStore, int32) {
				user, err := testtool.GenRdmUser()
				require.NoError(t, err)

				req := &service.UserCreateRequest{
					Password:  string(user.Password),
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Role:      string(user.Role),
					Email:     user.Email,
				}
				require.NotEmpty(t, req.RequestName())

				matcher, err := testtool.NewUserCreateReqMatcher(req)
				require.NoError(t, err)

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					CreateUser(gomock.Any(), matcher).
					Times(1).
					Return(user.ID, nil)

				return req, store, user.ID
			},
			TestFunc: func(req *service.UserCreateRequest, store *mock_model.MockStore, id int32) {
				srvc := service.NewService(store, Validate)
				n, err := srvc.User().Create(context.Background(), req)
				require.NoError(t, err)
				require.Equal(t, id, n)
			},
			CleanUpFunc: func() {
				Validate = nil
			},
		},
		{
			Name: "Role Validation Error",
			Init: func() {
				Validate = val.New()
				validator.RegisterValidator(
					Validate,
					validator.EnmusRole,
					validator.NewDefaultPasswordValidator(),
				)
			},
			SetupFunc: func(ctl *gomock.Controller) (*service.UserCreateRequest, *mock_model.MockStore, int32) {
				user, err := testtool.GenRdmUser()
				require.NoError(t, err)

				req := &service.UserCreateRequest{
					Password:  string(user.Password),
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Role:      rg.Must(rg.AlphaNum.GenRdmString(13)),
					Email:     user.Email,
				}
				require.NotEmpty(t, req.RequestName())

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)

				return req, store, user.ID
			},
			TestFunc: func(req *service.UserCreateRequest, store *mock_model.MockStore, id int32) {
				srvc := service.NewService(store, Validate)
				n, err := srvc.User().Create(context.Background(), req)
				var valErr val.ValidationErrors
				ok := errors.As(err, &valErr)
				require.True(t, ok)
				require.Error(t, err)
				require.Equal(t, int32(0), n)
			},
			CleanUpFunc: func() {
				Validate = nil
			},
		},
		{
			Name: "Password too long (bcrypt error)",
			Init: func() {
				Validate = val.New()
				leakyPWDValidator := validator.NewPasswordValidator(
					false, 0, 100, 0, 0, 0, 0,
				)
				validator.RegisterValidator(
					Validate,
					validator.EnmusRole,
					leakyPWDValidator,
				)
			},
			SetupFunc: func(ctl *gomock.Controller) (*service.UserCreateRequest, *mock_model.MockStore, int32) {
				user, err := testtool.GenRdmUser()
				require.NoError(t, err)

				req := &service.UserCreateRequest{
					Password:  rg.Must[string](rg.AlphaNum.GenRdmString(80)),
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Role:      string(user.Role),
					Email:     user.Email,
				}
				require.NotEmpty(t, req.RequestName())

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)

				return req, store, user.ID
			},
			TestFunc: func(req *service.UserCreateRequest, store *mock_model.MockStore, id int32) {
				srvc := service.NewService(store, Validate)
				n, err := srvc.User().Create(context.Background(), req)
				require.Error(t, err)
				require.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
				require.Equal(t, int32(0), n)
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

func TestGetUserAuthInfoPassword(t *testing.T) {
	type testCase struct {
		Name        string
		Init        func()
		SetupFunc   func(ctl *gomock.Controller) (string, *mock_model.MockStore, *model.GetUserAuthRow)
		TestFunc    func(email string, store *mock_model.MockStore, auth *model.GetUserAuthRow)
		CleanUpFunc func()
	}

	var Validate *val.Validate
	tcs := []testCase{
		{
			Name: "OK",
			Init: func() {
				Validate = val.New()
				validator.RegisterValidator(Validate, validator.NewDefaultPasswordValidator())
			},
			SetupFunc: func(ctl *gomock.Controller) (string, *mock_model.MockStore, *model.GetUserAuthRow) {
				user, err := testtool.GenRdmUser()
				require.NoError(t, err)

				row := &model.GetUserAuthRow{
					ID: user.ID, Email: user.Email, Password: user.Password,
				}
				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(row, nil)

				return user.Email, store, row
			},
			TestFunc: func(email string, store *mock_model.MockStore, auth *model.GetUserAuthRow) {
				srvc := service.NewService(store, Validate)
				row, err := srvc.User().GetAuthInfo(context.Background(), email)
				require.NoError(t, err)
				require.Equal(t, auth.ID, row.ID)
				require.Equal(t, auth.Email, row.Email)
				require.Equal(t, auth.Password, row.Password)
			},
			CleanUpFunc: func() {
				Validate = nil
			},
		},
		{
			Name: "Row with given email not found",
			Init: func() {
				Validate = val.New()
				validator.RegisterValidator(Validate, validator.NewDefaultPasswordValidator())
			},
			SetupFunc: func(ctl *gomock.Controller) (string, *mock_model.MockStore, *model.GetUserAuthRow) {
				user, err := testtool.GenRdmUser()
				require.NoError(t, err)

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(nil, pgx.ErrNoRows)

				return user.Email, store, nil
			},
			TestFunc: func(email string, store *mock_model.MockStore, auth *model.GetUserAuthRow) {
				srvc := service.NewService(store, Validate)
				row, err := srvc.User().GetAuthInfo(context.Background(), email)
				require.ErrorIs(t, err, pgx.ErrNoRows)
				require.Nil(t, row)
			},
			CleanUpFunc: func() {
				Validate = nil
			},
		},
		{
			Name: "Bad Email",
			Init: func() {
				Validate = val.New()
				validator.RegisterValidator(Validate, validator.NewDefaultPasswordValidator())
			},
			SetupFunc: func(ctl *gomock.Controller) (string, *mock_model.MockStore, *model.GetUserAuthRow) {
				email, err := rg.AlphaNum.GenRdmString(20)
				require.NoError(t, err)

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					GetUserAuth(gomock.Any(), gomock.Any()).
					Times(0)

				return email, store, nil
			},
			TestFunc: func(email string, store *mock_model.MockStore, auth *model.GetUserAuthRow) {
				srvc := service.NewService(store, Validate)
				row, err := srvc.User().GetAuthInfo(context.Background(), email)
				var valErr val.ValidationErrors
				ok := errors.As(err, &valErr)
				require.True(t, ok)
				require.Error(t, err)
				require.Nil(t, row)
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

func TestUpdateUserPasswordService(t *testing.T) {
	type testCase struct {
		Name        string
		Init        func()
		SetupFunc   func(ctl *gomock.Controller) (*service.UserUpdatePasswordRequest, *mock_model.MockStore, int64)
		TestFunc    func(req *service.UserUpdatePasswordRequest, store *mock_model.MockStore, nAffectedRow int64)
		CleanUpFunc func()
	}

	var Validate *val.Validate
	tcs := []testCase{
		{
			Name: "Create User OK",
			Init: func() {
				Validate = val.New()
				validator.RegisterValidator(
					Validate,
					validator.EnmusRole,
					validator.NewDefaultPasswordValidator(),
				)
			},
			SetupFunc: func(ctl *gomock.Controller) (*service.UserUpdatePasswordRequest, *mock_model.MockStore, int64) {
				req := &service.UserUpdatePasswordRequest{
					ID:       rand.Int31n(100_000) + 1,
					Password: string(rg.Must[[]byte](rg.GenRdmPwd(8, 30, 2, 2, 2, 2))),
				}
				require.NotEmpty(t, req.RequestName())

				matcher, err := testtool.NewUserUpdatePasswordReqMatcher(req)
				require.NoError(t, err)

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					UpdatePassword(gomock.Any(), matcher).
					Times(1).
					Return(int64(1), nil)

				return req, store, 1
			},
			TestFunc: func(req *service.UserUpdatePasswordRequest, store *mock_model.MockStore, nAffectedRow int64) {
				srvc := service.NewService(store, Validate)
				n, err := srvc.User().UpdatePassword(context.Background(), req)
				require.NoError(t, err)
				require.Equal(t, nAffectedRow, n)
			},
			CleanUpFunc: func() {
				Validate = nil
			},
		},
		{
			Name: "ID is invalid",
			Init: func() {
				Validate = val.New()
				validator.RegisterValidator(
					Validate,
					validator.EnmusRole,
					validator.NewDefaultPasswordValidator(),
				)
			},
			SetupFunc: func(ctl *gomock.Controller) (*service.UserUpdatePasswordRequest, *mock_model.MockStore, int64) {
				req := &service.UserUpdatePasswordRequest{
					ID:       -rand.Int31n(1_000),
					Password: string(rg.Must[[]byte](rg.GenRdmPwd(8, 30, 2, 2, 2, 2))),
				}
				require.NotEmpty(t, req.RequestName())

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					UpdatePassword(gomock.Any(), gomock.Any()).
					Times(0)

				return req, store, 0
			},
			TestFunc: func(req *service.UserUpdatePasswordRequest, store *mock_model.MockStore, nAffectedRow int64) {
				srvc := service.NewService(store, Validate)
				n, err := srvc.User().UpdatePassword(context.Background(), req)
				var valErr val.ValidationErrors
				ok := errors.As(err, &valErr)
				require.True(t, ok)
				require.Error(t, err)
				require.Equal(t, int64(0), n)
			},
			CleanUpFunc: func() {
				Validate = nil
			},
		},
		{
			Name: "Password too long (bcrypt error)",
			Init: func() {
				Validate = val.New()
				leakyPWDValidator := validator.NewPasswordValidator(
					false, 0, 100, 0, 0, 0, 0,
				)
				validator.RegisterValidator(
					Validate,
					validator.EnmusRole,
					leakyPWDValidator,
				)
			},
			SetupFunc: func(ctl *gomock.Controller) (*service.UserUpdatePasswordRequest, *mock_model.MockStore, int64) {
				req := &service.UserUpdatePasswordRequest{
					ID:       rand.Int31n(100_000),
					Password: string(rg.Must[[]byte](rg.GenRdmPwd(80, 100, 2, 2, 2, 2))),
				}
				require.NotEmpty(t, req.RequestName())

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					UpdatePassword(gomock.Any(), gomock.Any()).
					Times(0)

				return req, store, 0
			},
			TestFunc: func(req *service.UserUpdatePasswordRequest, store *mock_model.MockStore, nAffectedRow int64) {
				srvc := service.NewService(store, Validate)
				n, err := srvc.User().UpdatePassword(context.Background(), req)
				require.Error(t, err)
				require.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
				require.Equal(t, int64(0), n)
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

func TestDeleteUserService(t *testing.T) {
	type testCase struct {
		Name        string
		Init        func()
		SetupFunc   func(ctl *gomock.Controller) (int32, *mock_model.MockStore, int64)
		TestFunc    func(id int32, store *mock_model.MockStore, nAffectedRow int64)
		CleanUpFunc func()
	}

	var Validate *val.Validate
	tcs := []testCase{
		{
			Name: "Delete User OK",
			Init: func() {
				Validate = val.New()
				validator.RegisterValidator(
					Validate,
					validator.EnmusRole,
					validator.NewDefaultPasswordValidator(),
				)
			},
			SetupFunc: func(ctl *gomock.Controller) (int32, *mock_model.MockStore, int64) {
				id := rand.Int31n(100_000) + 1

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(id)).
					Times(1).
					Return(int64(1), nil)

				return id, store, 1
			},
			TestFunc: func(id int32, store *mock_model.MockStore, nAffectedRow int64) {
				srvc := service.NewService(store, Validate)
				n, err := srvc.User().Delete(context.Background(), id)
				require.NoError(t, err)
				require.Equal(t, nAffectedRow, n)
			},
			CleanUpFunc: func() {
				Validate = nil
			},
		},
		{
			Name: "ID is invalid",
			Init: func() {
				Validate = val.New()
				validator.RegisterValidator(
					Validate,
					validator.EnmusRole,
					validator.NewDefaultPasswordValidator(),
				)
			},
			SetupFunc: func(ctl *gomock.Controller) (int32, *mock_model.MockStore, int64) {
				id := -rand.Int31n(100_000)

				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(0)

				return id, store, 0
			},
			TestFunc: func(id int32, store *mock_model.MockStore, nAffectedRow int64) {
				srvc := service.NewService(store, Validate)
				n, err := srvc.User().Delete(context.Background(), id)
				var valErr val.ValidationErrors
				ok := errors.As(err, &valErr)
				require.True(t, ok)
				require.Error(t, err)
				require.Equal(t, int64(0), n)
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

func TestCleanUpUserService(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		nAffectedRow := rand.Int63n(1_000) + 1
		store := mock_model.NewMockStore(ctl)
		store.
			EXPECT().
			CleanUpUsers(gomock.Any()).
			Times(1).
			Return(nAffectedRow, nil)

		srvc := service.NewService(store, nil)
		n, err := srvc.User().CleanUp(context.Background())
		require.NoError(t, err)
		require.Equal(t, nAffectedRow, n)
	})

	t.Run("Time Out", func(t *testing.T) {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		db := make(chan struct{})
		defer close(db)

		nAffectedRow := rand.Int63n(1_000) + 1
		store := mock_model.NewMockStore(ctl)
		store.
			EXPECT().
			CleanUpUsers(gomock.Any()).
			Times(1).
			DoAndReturn(func(ctx context.Context) (n int64, err error) {
				select {
				case <-ctx.Done():
					return 0, errors.New("timeout")
				case <-db:
					return nAffectedRow, nil
				}
			})

		srvc := service.NewService(store, nil)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		go func(cancel context.CancelFunc) {
			time.Sleep(1 * time.Second)
			cancel()
		}(cancel)

		go func() {
			time.Sleep(3 * time.Second)
			db <- struct{}{}
		}()

		n, err := srvc.User().CleanUp(ctx)
		require.ErrorContains(t, err, "timeout")
		require.Equal(t, int64(0), n)
	})
}
