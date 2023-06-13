package convert

import (
	"math"

	"github.com/jackc/pgx/v5/pgtype"
)

type IntT interface {
	int | int8 | int16 | int32 | int64
}

type IntTo struct {
	i int64
}

func Must[T any](i T, err error) T {
	return i
}

func NewInTo[T IntT](i T) IntTo {
	return IntTo{int64(i)}
}

func (i IntTo) Int() (int, error) {
	i64, _ := i.Int64()
	if i64 > math.MaxInt || i64 < math.MinInt {
		return 0, ErrOverFlow
	}
	return int(i64), nil
}

func (i IntTo) Int64() (int64, error) {
	return i.i, nil
}

func (i IntTo) Int32() (int32, error) {
	i64, _ := i.Int64()
	if i64 > math.MaxInt32 || i64 < math.MinInt32 {
		return 0, ErrOverFlow
	}
	return int32(i.i), nil
}

func (i IntTo) Int16() (int16, error) {
	i64, _ := i.Int64()
	if i64 > math.MaxInt16 || i64 < math.MinInt16 {
		return 0, ErrOverFlow
	}
	return int16(i.i), nil
}

func (i IntTo) ToPgTypeInt4() (pgtype.Int4, error) {
	i32, err := i.Int32()
	if err != nil {
		return pgtype.Int4{Int32: 0, Valid: false}, err
	}
	return pgtype.Int4{Int32: i32, Valid: true}, nil
}

func (i IntTo) ToPgTypeInt2() (pgtype.Int2, error) {
	i16, err := i.Int16()
	if err != nil {
		return pgtype.Int2{Int16: 0, Valid: false}, err
	}
	return pgtype.Int2{Int16: i16, Valid: true}, nil
}
