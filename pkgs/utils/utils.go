package utils

import (
	"errors"
	"reflect"
	"strings"
)

func Unique[T comparable](v []T) []T {
	set := map[T]struct{}{}
	for _, e := range v {
		set[e] = struct{}{}
	}
	u := make([]T, 0, len(set))
	for k := range set {
		u = append(u, k)
	}
	return u
}

func GetStructTags(key string, s interface{}) ([]string, error) {
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		return nil, errors.New("bad type")
	}

	tags := []string{}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tag := strings.Split(f.Tag.Get(key), ",")[0]
		if tag != "" && tag != "-" {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

func GetFieldName(tag, key string, s interface{}) (fieldname string) {
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get(key), ",")[0]
		if v == tag {
			return f.Name
		}
	}
	return ""
}
