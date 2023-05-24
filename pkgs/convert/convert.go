package convert

import (
	"errors"
	"strconv"
)

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

func (s StrTo) Int() (int, error) {
	return strconv.Atoi(s.String())
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

var ErrOverFlow = errors.New("value overflow")

func (s StrTo) UInt32() (uint32, error) {
	v, err := s.Int()
	if err != nil {
		return uint32(0), err
	}

	if v > 4294967295 || v < 0 {
		return uint32(0), ErrOverFlow
	}
	return uint32(v), nil
}

func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}
