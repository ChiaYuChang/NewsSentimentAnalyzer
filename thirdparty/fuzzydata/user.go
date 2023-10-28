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
	model.RoleUser,
	model.RoleAdmin,
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

func (u UserItem) Password() ([]byte, error) {
	return bcrypt.GenerateFromPassword(u.RawPassword, bcrypt.DefaultCost)
}

func NewUsers(n int) User {
	n += 2
	us := User{}
	us.Item = make([]UserItem, n)

	var rts []time.Time
	rts = rg.GenRdnTimes(2, TIME_MIN, TIME_MAX)
	us.Item[0] = UserItem{
		Id:          uuid.New(),
		RawPassword: []byte("password"),
		FirstName:   rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
		LastName:    rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
		Role:        UserRole[0],
		Email:       "test@example.com",
		CreatedAt:   rts[0],
		UpdatedAt:   rts[1],
	}

	rts = rg.GenRdnTimes(2, TIME_MIN, TIME_MAX)
	us.Item[1] = UserItem{
		Id:          uuid.New(),
		RawPassword: []byte("password"),
		FirstName:   rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
		LastName:    rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
		Role:        UserRole[1],
		Email:       "admin@example.com",
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
