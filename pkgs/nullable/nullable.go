package nullable

import (
	"bytes"
	"log"
	"strconv"
	"time"
)

type String[T ~string] struct {
	Value T
	Valid bool
}

func (s *String[T]) UnmarshalJSON(bs []byte) error {
	log.Println("using NullableString.UnmarshalJSON")
	if bytes.Equal(bs, []byte("null")) {
		(*s) = String[T]{}
	} else {
		(*s).Valid = true
		(*s).Value = T(string(bytes.Trim(bytes.Trim(bs, "\""), "'")))
	}
	return nil
}

func (s *String[T]) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return []byte("\"" + s.Value + "\""), nil
	}
	return []byte("null"), nil
}

type Integer[T ~int] struct {
	Value T
	Valid bool
}

func (i *Integer[T]) UnmarshalJSON(bs []byte) error {
	if bytes.Equal(bs, []byte("null")) {
		(*i).Valid = false
	} else {
		(*i).Valid = true
		j, err := strconv.Atoi(string(bs))
		if err != nil {
			return err
		}
		(*i).Value = T(j)
	}
	return nil
}

func (i *Integer[T]) MarshalJSON() ([]byte, error) {
	if i.Valid {
		si := strconv.Itoa(int(i.Value))
		return []byte(si), nil
	}
	return []byte("null"), nil
}

type Float32[T ~float32] struct {
	Value T
	Valid bool
}

func (f32 *Float32[T]) UnmarshalJSON(bs []byte) error {
	if bytes.Equal(bs, []byte("null")) {
		(*f32).Valid = false
	} else {
		(*f32).Valid = true
		f, err := strconv.ParseFloat(string(bs), 32)
		if err != nil {
			return err
		}
		(*f32).Value = T(float32(f))
	}
	return nil
}

func (f32 *Float32[T]) MarshalJSON() ([]byte, error) {
	if f32.Valid {
		return []byte(strconv.FormatFloat(float64(f32.Value), 'f', -1, 32)), nil
	}
	return []byte("null"), nil
}

type Float64[T ~float64] struct {
	Value T
	Valid bool
}

func (f64 *Float64[T]) UnmarshalJSON(bs []byte) error {
	if bytes.Equal(bs, []byte("null")) {
		(*f64).Valid = false
	} else {
		(*f64).Valid = true
		f, err := strconv.ParseFloat(string(bs), 64)
		if err != nil {
			return err
		}
		(*f64).Value = T(float32(f))
	}
	return nil
}

func (f64 *Float64[T]) MarshalJSON() ([]byte, error) {
	if f64.Valid {
		return []byte(strconv.FormatFloat(float64(f64.Value), 'f', -1, 64)), nil
	}
	return []byte("null"), nil
}

type Time struct {
	Value  time.Time
	Valid  bool
	Layout string
}

func (i *Time) UnmarshalJSON(bs []byte) error {
	if bytes.Equal(bs, []byte("null")) {
		(*i).Valid = false
	} else {
		t, err := time.Parse(i.Layout, string(bs))
		if err != nil {
			return err
		}
		(*i).Valid = true
		(*i).Value = t
	}
	return nil
}

func (i *Time) MarshalJSON() ([]byte, error) {
	if i.Valid {
		si := i.Value.Format(i.Layout)
		return []byte(si), nil
	}
	return []byte("null"), nil
}
