package service

import (
	"fmt"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	"golang.org/x/net/context"
)

func (srvc userService) Service() Service {
	return Service(srvc)
}

type UserCreateRequest struct {
	Password  []byte `validate:"required,min=8,max=30"`
	FirstName string `validate:"required,min=1,max=20"`
	LastName  string `validate:"required,min=1,max=20"`
	Role      string `validate:"required,role"`
	Email     string `validate:"required,email"`
}

func (r UserCreateRequest) RequestName() string {
	return "user-create-req"
}

type UserGetAuthInfoRequest struct {
	Email string `validate:"email"`
}

func (r UserGetAuthInfoRequest) RequestName() string {
	return "user-get-auth-request"
}

type UserUpdatePasswordRequest struct {
	ID       int32  `validate:"required"`
	Password []byte `validate:"required,password"`
}

func (r UserUpdatePasswordRequest) RequestName() string {
	return "user-update-password-request"
}

type UserDeleteRequest struct {
	ID int32 `validate:"required,min=1"`
}

func (r UserDeleteRequest) RequestName() string {
	return "user-delete-req"
}

func (srvc userService) Create(ctx context.Context, req *UserCreateRequest) error {
	var err error
	err = srvc.Struct(req)
	if err != nil {
		return err
	}

	err = srvc.Store.CreateUser(ctx, &model.CreateUserParams{
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      model.Role(req.Role),
		Email:     req.Email,
	})
	return err
}

func (srvc userService) GetAuthInfo(ctx context.Context, req *UserGetAuthInfoRequest) (*model.GetUserAuthRow, error) {
	if err := srvc.Struct(req); err != nil {
		return nil, fmt.Errorf("Validation error: %w", err)
	}
	return srvc.Store.GetUserAuth(ctx, req.Email)
}

func (srvc userService) UpdatePassword(ctx context.Context, req *UserUpdatePasswordRequest) error {
	if err := srvc.Validate.Struct(req); err != nil {
		return fmt.Errorf("Validation error: %w", err)
	}
	return srvc.Store.UpdatePassword(ctx,
		&model.UpdatePasswordParams{
			ID:       req.ID,
			Password: req.Password,
		})
}

func (srvc userService) Delete(ctx context.Context, req *UserDeleteRequest) error {
	if err := srvc.Validate.Struct(req); err != nil {
		return fmt.Errorf("Validation error: %w", err)
	}
	return srvc.Store.DeleteUser(ctx, req.ID)
}
