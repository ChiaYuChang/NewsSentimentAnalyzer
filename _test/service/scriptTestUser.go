package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model/testtool"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type CreateUpateThenDeleteUser struct {
	n int
}

func (scrpt CreateUpateThenDeleteUser) Do(
	ctx context.Context, srvc service.Service, logger Logger) Job {
	var err error

	jid := ctx.Value("jid").(int)

	logger.Print("\t- Check Create User...")
	userOri, _ := testtool.GenRdmUser()
	userOri.Password = []byte("P@ssword123")
	userCreateReq := &service.UserCreateRequest{
		Password:  string(userOri.Password),
		FirstName: userOri.FirstName,
		LastName:  userOri.LastName,
		Role:      string(userOri.Role),
		Email:     userOri.Email,
	}

	userOri.ID, err = srvc.User().Create(context.Background(), userCreateReq)
	if err != nil {
		logger.Println("Failed")
		return Job{jid, err}
	}
	logger.Println("OK")

	logger.Print("\t- Check Get User Auth...")
	userOriAuth, err := srvc.User().GetAuthInfo(context.Background(), userOri.Email)
	if err != nil {
		logger.Println("Failed")
		return Job{jid, err}
	}
	if err := bcrypt.CompareHashAndPassword(
		userOriAuth.Password, userOri.Password); err != nil {
		return Job{jid, err}
	}
	logger.Println("OK")

	logger.Print("\t- Check Update User Password...")
	userUpdated := testtool.CloneUser(userOri)
	userUpdated.Password, _ = rg.GenRdmPwd(10, 20, 1, 1, 1, 1)
	userUpdatePasswordReq := &service.UserUpdatePasswordRequest{
		ID:       userUpdated.ID,
		Password: string(userUpdated.Password),
	}

	n, err := srvc.User().UpdatePassword(context.Background(), userUpdatePasswordReq)
	if n != 1 {
		logger.Println("Failed")
		return Job{jid, fmt.Errorf("%w: expected: %d, actual: %d",
			ErrorAffectMoreThanExpectedRows, 1, n)}
	}
	if err != nil {
		logger.Println("Failed")
		return Job{jid, err}
	}

	userUpdatedAuth, err := srvc.User().GetAuthInfo(context.Background(), userUpdated.Email)
	if err != nil {
		return Job{jid, err}
	}

	if err := bcrypt.CompareHashAndPassword(
		userUpdatedAuth.Password, userUpdated.Password); err != nil {
		logger.Println("Failed")
		return Job{jid, err}
	}
	logger.Println("OK")

	logger.Print("\t- Check Delete User...")
	n, err = srvc.User().Delete(context.Background(), userOri.ID)
	if n != 1 {
		return Job{jid, fmt.Errorf("%w: expected: %d, actual: %d",
			ErrorAffectMoreThanExpectedRows, 1, n)}
	}
	if err != nil {
		logger.Println("Failed")
		return Job{jid, err}
	}
	_, err = srvc.User().GetAuthInfo(context.Background(), userOri.Email)
	if !errors.Is(err, pgx.ErrNoRows) {
		return Job{jid, err}
	}

	logger.Println("OK")
	return Job{jid, nil}
}

func (scrpt CreateUpateThenDeleteUser) Steps() []string {
	return []string{
		"Create user",
		"Get User Auth",
		"Update user password",
		"Delete user",
	}
}

func (scrpt CreateUpateThenDeleteUser) Description() string {
	return strings.Join(scrpt.Steps(), " -> ")
}

func (scrpt CreateUpateThenDeleteUser) N() int {
	return scrpt.n
}
