package pageform

import (
	"strings"
	"time"
)

type SearchIn struct {
	InTitle       bool `form:"in-title"`
	InDescription bool `form:"in-description"`
	InContent     bool `form:"in-content"`
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

type TimeRange struct {
	Form time.Time `form:"from-time" validate:"lte"`
	To   time.Time `form:"to-time" validate:"lte"`
}
