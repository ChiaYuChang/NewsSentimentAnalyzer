package service

import (
	"fmt"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

func (srvc userService) Service() Service {
	return Service(srvc)
}

func (srvc userService) Create(
	ctx context.Context, req *UserCreateRequest) (id int32, err error) {
	if err := srvc.validate.Struct(req); err != nil {
		return 0, fmt.Errorf("validation error: %w", err)
	}

	if params, err := req.ToParams(); err != nil {
		return 0, fmt.Errorf("error while ToParams error: %w", err)
	} else {
		return srvc.store.CreateUser(ctx, params)
	}
}

func (srvc userService) GetAuthInfo(
	ctx context.Context, email string) (*model.GetUserAuthRow, error) {
	if err := srvc.validate.Var(email, "required,email"); err != nil {
		return nil, fmt.Errorf("Validation error: %w", err)
	}
	return srvc.store.GetUserAuth(ctx, email)
}

func (srvc userService) UpdatePassword(
	ctx context.Context, req *UserUpdatePasswordRequest) (n int64, err error) {
	if err := srvc.validate.Struct(req); err != nil {
		return 0, err
	}

	if params, err := req.ToParams(); err != nil {
		return 0, fmt.Errorf("error while ToParams error: %w", err)
	} else {
		return srvc.store.UpdatePassword(ctx, params)
	}
}

func (srvc userService) Delete(
	ctx context.Context, Id int32) (n int64, err error) {
	if err := srvc.validate.Var(Id, "required,min=1"); err != nil {
		return 0, err
	}
	return srvc.store.DeleteUser(ctx, Id)
}

func (srvc userService) CleanUp(ctx context.Context) (n int64, err error) {
	return srvc.store.CleanUpUsers(ctx)
}

type UserCreateRequest struct {
	Password  string `validate:"required,password"`
	FirstName string `validate:"required,min=1,max=20"`
	LastName  string `validate:"required,min=1,max=20"`
	Role      string `validate:"required,role"`
	Email     string `validate:"required,email"`
}

func (req UserCreateRequest) RequestName() string {
	return "user-create-req"
}

func (req UserCreateRequest) ToParams() (*model.CreateUserParams, error) {
	if ePwd, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	); err != nil {
		return nil, fmt.Errorf("error while bcryot,genpwd: %w", err)
	} else {
		return &model.CreateUserParams{
			Password:  ePwd,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Role:      model.Role(req.Role),
			Email:     req.Email,
		}, nil
	}
}

type UserUpdatePasswordRequest struct {
	ID       int32  `validate:"required,min=1"`
	Password string `validate:"required,password"`
}

func (r UserUpdatePasswordRequest) RequestName() string {
	return "user-update-password-request"
}

func (req UserUpdatePasswordRequest) ToParams() (*model.UpdatePasswordParams, error) {
	if ePwd, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	); err != nil {
		return nil, fmt.Errorf("error while bcryot,genpwd: %w", err)
	} else {
		return &model.UpdatePasswordParams{
			ID:       req.ID,
			Password: ePwd,
		}, nil
	}
}
