package pageform

import (
	"fmt"
	"strings"
	"time"

	val "github.com/go-playground/validator/v10"
)

var LocationValidator = locationValidator{}

type locationValidator struct{}

func (tzv locationValidator) Tag() string {
	return "loc"
}

func (tzv locationValidator) ValFunc() val.Func {
	return func(fl val.FieldLevel) bool {
		tz := fl.Field().Interface().(string)
		if tz != "" {
			_, err := time.LoadLocation(tz)
			if err != nil {
				return false
			}
		}
		return true
	}
}

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

	if len(ss) == 3 || len(ss) == 0 {
		return ""
	}
	return strings.Join(ss, ",")
}

type TimeRange struct {
	Form     time.Time `form:"from-time"`
	To       time.Time `form:"to-time"`
	Location string    `form:"timezone"  validate:"required,loc"`
}

func (tr TimeRange) String() string {
	return tr.ToString("")
}

func (tr TimeRange) ToString(prefix string) string {
	sb := strings.Builder{}
	if !tr.Form.IsZero() {
		sb.WriteString(
			fmt.Sprintf("%s- From    : %s (%s)\n",
				prefix,
				tr.Form.Format(time.DateOnly),
				tr.Location))
	}
	if !tr.To.IsZero() {
		sb.WriteString(
			fmt.Sprintf("%s- To      : %s (%s)\n",
				prefix,
				tr.To.Format(time.DateOnly),
				tr.Location))
	}
	return sb.String()
}

func (tr *TimeRange) ToUTC() *TimeRange {
	loc, err := time.LoadLocation(tr.Location)
	if err != nil {
		loc = time.UTC
	}

	tr.Form = time.Date(tr.Form.Year(), tr.Form.Month(), tr.Form.Day(), 23, 59, 59, 0, loc).UTC()
	tr.To = time.Date(tr.To.Year(), tr.To.Month(), tr.To.Day(), 0, 0, 0, 0, loc).UTC()
	return tr
}
