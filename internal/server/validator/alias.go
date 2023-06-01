package validator

import (
	"strings"

	val "github.com/go-playground/validator/v10"
)

type Alias struct {
	alias string
	tags  []string
}

func (a Alias) Tags() string {
	return strings.Join(a.tags, ",")
}

func RegisterAlias(val *val.Validate, aliases ...Alias) {
	for _, alias := range aliases {
		val.RegisterAlias(alias.alias, alias.Tags())
	}
}

var APIName = Alias{
	alias: "api_name",
	tags:  []string{"min=3", "max=20"},
}

var QueryLimit = Alias{
	alias: "query_limit",
	tags:  []string{"min=1", "max=100"},
}
