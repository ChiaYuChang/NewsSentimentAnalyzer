package pageform

import (
	"fmt"
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

	if len(ss) == 3 {
		return ""
	}
	return strings.Join(ss, ",")
}

type TimeRange struct {
	Form time.Time `form:"from-time" validate:"lte"`
	To   time.Time `form:"to-time" validate:"lte"`
}

func (tr TimeRange) String() string {
	return tr.ToString("")
}

func (tr TimeRange) ToString(prefix string) string {
	sb := strings.Builder{}
	if !tr.Form.IsZero() {
		sb.WriteString(fmt.Sprintf("%s- From    : %s\n", prefix, tr.Form.Format(time.DateOnly)))
	}
	if !tr.To.IsZero() {
		sb.WriteString(fmt.Sprintf("%s- To      : %s\n", prefix, tr.To.Format(time.DateOnly)))
	}
	return sb.String()
}
