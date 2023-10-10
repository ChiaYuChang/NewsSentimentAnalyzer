package parser

import (
	"errors"
	"fmt"
	"strings"
)

type Errors map[string]error

func (e Errors) Error() string {
	sb := strings.Builder{}
	sb.WriteString("multiple errors:\n")
	for k, v := range e {
		sb.WriteString(fmt.Sprintf("%-10s: %s\n", k, v.Error()))
	}
	return sb.String()
}

func (e Errors) Add(key string, err error) {
	e[key] = err
}

func (e Errors) HasError(err error) bool {
	for _, v := range e {
		if errors.Is(v, err) {
			return true
		}
	}
	return false
}
