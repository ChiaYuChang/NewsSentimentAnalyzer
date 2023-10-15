package newsdata_test

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	newsdata "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/NEWSDATA"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/go-playground/form"
	mod "github.com/go-playground/mold/v4/modifiers"
	scr "github.com/go-playground/mold/v4/scrubbers"
	val "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestNEWSDATAIOFormValidationStruct(t *testing.T) {
	val := val.New()
	err := validator.RegisterValidator(
		val,
		newsdata.LanguageValidator,
		newsdata.CountryValidator,
		newsdata.CategoryValidator,
	)
	require.NoError(t, err)

	type valCatStruct struct {
		Category string `validate:"newsdata_cat"`
	}

	type valCtryStruct struct {
		Country string `validate:"newsdata_ctry"`
	}

	type valLangStruct struct {
		Language string `validate:"newsdata_lang"`
	}

	require.NoError(t, val.Var("business", newsdata.CategoryValidator.Tag()))
	require.NoError(t, val.Struct(valCatStruct{Category: "business"}))
	require.Error(t, val.Var("xx", newsdata.CategoryValidator.Tag()))
	require.Error(t, val.Struct(valCatStruct{Category: "xx"}))

	require.NoError(t, val.Var("tw", newsdata.CountryValidator.Tag()))
	require.NoError(t, val.Struct(valCtryStruct{Country: "tw"}))
	require.Error(t, val.Var("xx", newsdata.CountryValidator.Tag()))
	require.Error(t, val.Struct(valCtryStruct{Country: "xx"}))

	require.NoError(t, val.Var("en", newsdata.LanguageValidator.Tag()))
	require.NoError(t, val.Struct(valLangStruct{Language: "en"}))
	require.Error(t, val.Var("xx", newsdata.LanguageValidator.Tag()))
	require.Error(t, val.Struct(valLangStruct{Language: "xx"}))
}

type User struct {
	Name     string    `mod:"trim"       form:"name"  validate:"required"       scrub:"name"`
	Email    string    `mod:"trim,lcase" form:"email" validate:"required,email" scrub:"emails"`
	Birthday time.Time `mod:"trim"       form:"birthday"`
}

func (f User) String() string {
	sb := strings.Builder{}
	sb.WriteString("User:\n")
	sb.WriteString(fmt.Sprintf("\t- Name     : %s\n", f.Name))
	sb.WriteString(fmt.Sprintf("\t- Email    : %s\n", f.Email))
	sb.WriteString(fmt.Sprintf("\t- Birthday : %s\n", f.Birthday))
	return sb.String()
}

func parseForm() url.Values {
	return url.Values{
		"name":     []string{"  joeybloggs  "},
		"email":    []string{"Dean.Karn@gmail.com  "},
		"birthday": []string{""},
	}
}

func TestMod(t *testing.T) {
	modifier := mod.New()
	scrubber := scr.New()
	decoder := form.NewDecoder()
	decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		val := strings.TrimSpace(vals[0])
		if len(val) == 0 {
			return time.Time{}, nil
		}
		return time.Parse(time.DateOnly, val)
	}, time.Time{})

	values := parseForm()

	var user User
	err := decoder.Decode(&user, values)
	require.NoError(t, err)
	t.Log(user)

	modifier.Struct(context.Background(), &user)
	t.Log(user)

	scrubber.Struct(context.Background(), &user)
	t.Log(user)

}
