package main

import (
	"crypto/rand"
	"errors"
	"math"
	"unicode/utf8"
)

type CharSet struct {
	setSize  uint8
	setChars []rune
	setBytes [][]byte
	setMap   map[rune]bool
	Check    bool
}

var (
	Digit, _         = NewCharSetFromStr("0123456789")
	AlphabetLower, _ = NewCharSetFromStr("abcdefghijklmnopqrstuvwxyz")
	AlphabetUpper, _ = NewCharSetFromStr("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	Alphabet, _      = AlphabetUpper.Merge(AlphabetLower)
	AlphaNum, _      = Alphabet.Merge(Digit)
	Special, _       = NewCharSetFromStr("!@%$%^&*")
	Password, _      = AlphaNum.Merge(Special)
	Punct, _         = NewCharSetFromStr("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")
	Blank, _         = NewCharSetFromStr(" \t")
)

var ErrCharSetTooLarge = errors.New("the size of a CharSet should be less than 256")

func NewCharSetFromStr(st string) (CharSet, error) {
	return NewCharSet([]rune(st))
}

func NewCharSet(rn []rune) (CharSet, error) {
	set := map[rune]struct{}{}
	for _, r := range rn {
		set[r] = struct{}{}
	}

	if len(set) > math.MaxUint8 {
		return *new(CharSet), ErrCharSetTooLarge
	}

	noDuplicatedRunes := make([]rune, 0, len(set))
	for r := range set {
		noDuplicatedRunes = append(noDuplicatedRunes, r)
	}

	cs := CharSet{
		setSize:  uint8(len(set)),
		setChars: noDuplicatedRunes,
		setMap:   nil,
		Check:    false,
	}
	cs.setBytes = cs.genByteSliceFromRunes(rn)
	return cs, nil
}

func (c CharSet) genByteSliceFromRunes(rs []rune) [][]byte {
	bs := make([][]byte, len(rs))
	for i, r := range rs {
		b := make([]byte, utf8.RuneLen(r))
		utf8.EncodeRune(b, r)
		bs[i] = b
	}
	return bs
}

func (c CharSet) SetSize() uint8 {
	return c.setSize
}

func (c CharSet) GetRune(i uint8) rune {
	return c.setChars[i%c.setSize]
}

func (c CharSet) GetByte(i uint8) []byte {
	return c.setBytes[i%c.setSize]
}

func (c1 CharSet) Merge(c2 CharSet) (CharSet, error) {
	setMap := make(map[rune]bool, c1.SetSize()+c2.SetSize())
	size := int64(0)
	for _, r := range c1.setChars {
		_, ok := setMap[r]
		if !ok {
			setMap[r] = true
			size += 1
		}
	}

	for _, r := range c2.setChars {
		_, ok := setMap[r]
		if !ok {
			setMap[r] = true
			size += 1
		}
	}

	setChars := make([]rune, size)
	i := 0
	for k := range setMap {
		setChars[i] = k
		i += 1
	}

	return NewCharSet(setChars)
}

func (c CharSet) GenerateRandomString(length int) (string, error) {
	index := make([]uint8, length)

	if _, err := rand.Read(index); err != nil {
		return "", err
	}

	randomRunes := make([]rune, length)
	for i, j := range index {
		randomRunes[i] = c.GetRune(j)
	}
	return string(randomRunes), nil
}
