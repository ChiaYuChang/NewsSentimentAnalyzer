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

func CloneUser(u *model.User) *model.User {
	pwd := make([]byte, len(u.Password))
	copy(pwd, u.Password)
	return &model.User{
		ID:        u.ID,
		Password:  pwd,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		Email:     u.Email,
	}
}

func GenRdmAPI() (*model.Api, error) {
	name, err := rg.Alphabet.GenRdmString(rand.Intn(18) + 3)
	if err != nil {
		return nil, err
	}

	apiType := model.ApiTypeSource
	if rand.Float64() > 0.8 {
		apiType = model.ApiTypeLanguageModel
	}

	return &model.Api{
		ID:   int16(rand.Intn(10_000)),
		Name: name,
		Type: apiType,
	}, nil
}

// generate a random apikey model. if 'owner', and 'api_id' <= 0, they will
// be generated randomly.
func GenRdmAPIKey(owner int32, api_id int16) (*model.Apikey, error) {
	if owner <= 0 {
		owner = rand.Int31n(1_000_000)
	}

	if api_id <= 0 {
		api_id = int16(rand.Intn(1_000))
	}

	keyLen := rand.Intn(30) + 30
	if key, err := rg.Password.GenRdmString(keyLen); err != nil {
		return nil, err
	} else {
		return &model.Apikey{
			ID:    rand.Int31n(100_000_000),
			Owner: owner,
			ApiID: api_id,
			Key:   key,
		}, nil
	}
}

func CloneAPI(a *model.Api) *model.Api {
	return &model.Api{
		ID:        a.ID,
		Name:      a.Name,
		Type:      a.Type,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		DeletedAt: a.DeletedAt,
	}
}
