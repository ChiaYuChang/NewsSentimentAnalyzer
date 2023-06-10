package pageform

import (
	"strings"
	"time"
)

type AuthInfo struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type SignUpInfo struct {
	AuthInfo
	FirstName string `form:"first-name"`
	LastName  string `form:"last-name"`
}

type GNewsHeadlines struct {
	Keyword  string    `form:"keyword"`
	Language string    `form:"language"`
	Country  string    `form:"country"`
	Category string    `form:"category"`
	Form     time.Time `form:"from-time"`
	To       time.Time `form:"to-time"`
}

type SearchIn struct {
	InTitle       bool `form:"title"`
	InDescription bool `form:"description"`
	InContent     bool `form:"content"`
}

func (s SearchIn) String() string {
	ss := make([]string, 0, 3)
	if s.InTitle {
		ss = append(ss, "title")
	}

	if s.InDescription {
		ss = append(ss, "description")
	}

	if s.InContent {
		ss = append(ss, "content")
	}
	return strings.Join(ss, ",")
}

type GNewsSearch struct {
	SearchIn
	Keyword  string    `form:"keyword"`
	Language string    `form:"language"`
	Country  string    `form:"country"`
	From     time.Time `form:"from-time"`
	To       time.Time `form:"to-time"`
}
