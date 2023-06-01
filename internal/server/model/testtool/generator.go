package testtool

import (
	"math/rand"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
)

func GenRdmUser() (*model.User, error) {
	pwd, err := rg.GenRdmPwd(8, 30, 1, 1, 1, 1)
	if err != nil {
		return nil, err
	}

	fnm, err := rg.Alphabet.GenRdmString(10)
	if err != nil {
		return nil, err
	}

	lnm, err := rg.Alphabet.GenRdmString(13)
	if err != nil {
		return nil, err
	}

	email, err := rg.GenRdmEmail(rg.AlphaNum, rg.Alphabet)
	if err != nil {
		return nil, err
	}

	role := model.RoleAdmin
	if rand.Float64() > 0.5 {
		role = model.RoleUser
	}

	return &model.User{
		ID:        rand.Int31() + 1,
		Password:  []byte(pwd),
		FirstName: fnm,
		LastName:  lnm,
		Role:      role,
		Email:     email,
	}, nil
}
