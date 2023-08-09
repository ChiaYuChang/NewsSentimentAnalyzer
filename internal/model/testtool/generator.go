package testtool

import (
	"math/rand"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/convert"
	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
		ID:        uuid.New(),
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

func GenRdmAPI(id int16) (*model.Api, error) {
	if id < 1 {
		id = int16(rand.Intn(10_000))
	}

	name, err := rg.Alphabet.GenRdmString(rand.Intn(18) + 3)
	if err != nil {
		return nil, err
	}

	apiType := model.ApiTypeSource
	if rand.Float64() > 0.8 {
		apiType = model.ApiTypeLanguageModel
	}

	return &model.Api{
		ID:   id,
		Name: name,
		Type: apiType,
	}, nil
}

// generate a random apikey model. if 'owner', and 'api_id' <= 0, they will
// be generated randomly.
func GenRdmAPIKey(owner uuid.UUID, api_id int16) (*model.Apikey, error) {
	if owner == uuid.Nil {
		owner = uuid.New()
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

func GenRdmAPIEndpoints(epId int32, apiId int16) (*model.Endpoint, error) {
	if epId < 1 {
		epId = rand.Int31n(10_000) + 1
	}

	if apiId < 1 {
		apiId = int16(rand.Int31n(100) + 1)
	}

	now := time.Now()
	ep := &model.Endpoint{
		ID:           epId,
		Name:         rg.Must[string](rg.AlphaNum.GenRdmString(rand.Intn(25) + 5)),
		ApiID:        apiId,
		TemplateName: rg.Must[string](rg.AlphaNum.GenRdmString(rand.Intn(20)+5)) + ".gotmpl",
		CreatedAt:    convert.TimeTo(now).ToPgTimeStampZ(),
		UpdatedAt:    convert.TimeTo(now.Add(time.Duration(rand.Intn(24*10)+1) * time.Hour)).ToPgTimeStampZ(),
		DeletedAt:    pgtype.Timestamptz{Valid: false},
	}

	return ep, nil
}

func CloneEndpoint(ep *model.Endpoint) *model.Endpoint {
	return &model.Endpoint{
		ID:           ep.ID,
		Name:         ep.Name,
		ApiID:        ep.ApiID,
		TemplateName: ep.TemplateName,
		CreatedAt:    ep.CreatedAt,
		UpdatedAt:    ep.UpdatedAt,
		DeletedAt:    ep.DeletedAt,
	}
}
