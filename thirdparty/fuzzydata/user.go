package main

import (
	"math/rand"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var UserRole = []model.Role{
	RoleUser,
	RoleAdmin,
}

type Password []byte

func (pws Password) String() string {
	return string(pws)
}

type User struct {
	Item []UserItem
	N    int
}

type UserItem struct {
	Id          uuid.UUID  `json:"id"`
	RawPassword []byte     `json:"raw_password"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Role        model.Role `json:"role"`
	Email       string     `json:"email"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (u UserItem) Password() (string, error) {
	b, err := bcrypt.GenerateFromPassword(u.RawPassword, bcrypt.DefaultCost)
	return string(b), err
}

func NewUsers(n int) User {
	n += 2
	us := User{}
	us.Item = make([]UserItem, n)

	var rts []time.Time
	rts = rg.GenRdnTimes(2, TIME_MIN, TIME_MAX)
	us.Item[0] = UserItem{
		Id:          TEST_USER_UID,
		RawPassword: []byte(TEST_USER_PASSWORD),
		FirstName:   rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
		LastName:    rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
		Role:        RoleUser,
		Email:       TEST_USER_EAMIL,
		CreatedAt:   rts[0],
		UpdatedAt:   rts[1],
	}

	rts = rg.GenRdnTimes(2, TIME_MIN, TIME_MAX)
	us.Item[1] = UserItem{
		Id:          TEST_ADMIN_USER_UID,
		RawPassword: []byte(TEST_USER_PASSWORD),
		FirstName:   rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
		LastName:    rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
		Role:        RoleAdmin,
		Email:       TEST_ADMIN_USER_EAMIL,
		CreatedAt:   rts[0],
		UpdatedAt:   rts[1],
	}

	rs := NewSampler(UserRole, []float64{0.99, 0.01})
	for i := 2; i < n; i++ {
		rawPwd, _ := rg.GenRdmPwd(8, 30, 1, 1, 1, 1)
		rts := rg.GenRdnTimes(2, TIME_MIN, TIME_MAX)
		u := UserItem{
			Id:          uuid.New(),
			RawPassword: rawPwd,
			FirstName:   rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
			LastName:    rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
			Role:        rs.Get(),
			Email:       rg.Must(rg.GenRdmEmail(rg.AlphaNum, rg.AlphabetLower)),
			CreatedAt:   rts[0],
			UpdatedAt:   rts[1],
		}
		us.Item[i] = u
	}

	us.N = n + 1
	return us
}
