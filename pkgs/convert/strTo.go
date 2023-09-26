package convert

import (
	"crypto/md5"
	"crypto/sha1"
	"math"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

func (s StrTo) Bytes() []byte {
	return []byte(s.String())
}

func (s StrTo) MD5Hash() ([16]byte, error) {
	md5Hasher := md5.New()
	_, err := md5Hasher.Write(s.Bytes())
	if err != nil {
		return [16]byte{}, err
	}
	return md5.Sum(nil), nil
}

func (s StrTo) SH1Hash() ([]byte, error) {
	sh1Hasher := sha1.New()
	_, err := sh1Hasher.Write(s.Bytes())
	if err != nil {
		return nil, err
	}
	return sh1Hasher.Sum(nil), nil
}

func (s StrTo) Int() (int, error) {
	return strconv.Atoi(s.String())
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

func (s StrTo) Int32() (int32, error) {
	if i, err := s.Int(); err != nil {
		return 0, err
	} else {
		return int32(i), nil
	}
}

func (s StrTo) MustInt32() int32 {
	v, _ := s.Int32()
	return v
}

func (s StrTo) UInt32() (uint32, error) {
	v, err := s.Int()
	if err != nil {
		return uint32(0), err
	}

	if v > math.MaxUint32 || v < 0 {
		return uint32(0), ErrOverFlow
	}
	return uint32(v), nil
}

func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}

func (s StrTo) PgText() pgtype.Text {
	return pgtype.Text{
		String: s.String(),
		Valid:  len(s) > 0,
	}
}
