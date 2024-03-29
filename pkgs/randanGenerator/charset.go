package randangenerator

import (
	"bytes"
	crand "crypto/rand"
	"errors"
	"math"
	"math/big"
	mrand "math/rand"
	"unicode/utf8"
)

type Rand interface {
	GetIntn(n int) int
	Read(dest []uint8) (n int, err error)
}

type CRand struct{}

func (c *CRand) GetIntn(n int) int {
	i, _ := crand.Int(crand.Reader, big.NewInt(int64(n)))
	return int(i.Int64())
}

func (c *CRand) Read(dest []uint8) (n int, err error) {
	return crand.Read(dest)
}

type MRand struct {
	*mrand.Rand
}

func NewMRand(seed int64) *MRand {
	r := &MRand{mrand.New(mrand.NewSource(seed))}
	return r
}

func (m *MRand) GetIntn(n int) int {
	return m.Intn(n)
}

func (m *MRand) Read(dest []uint8) (n int, err error) {
	for i := range dest {
		dest[i] = uint8(m.GetIntn(math.MaxUint8))
	}
	return len(dest), nil
}

type CharSet struct {
	setSize         uint8
	setChars        []rune
	setBytes        [][]byte
	setMap          map[rune]bool
	IsAllSingleByte bool
	Rand            Rand
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
	}
	cs.setBytes, cs.IsAllSingleByte = cs.genByteSliceFromRunes(rn)
	return cs, nil
}

func (c CharSet) genByteSliceFromRunes(rs []rune) ([][]byte, bool) {
	bs := make([][]byte, len(rs))
	isAllSingleByte := true

	for i, r := range rs {
		b := make([]byte, utf8.RuneLen(r))
		utf8.EncodeRune(b, r)
		bs[i] = b
		if len(b) > 1 {
			isAllSingleByte = false
		}
	}
	return bs, isAllSingleByte
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

func (c *CharSet) SetRand(rand Rand) CharSet {
	c.Rand = rand
	return *c
}

func (c1 CharSet) Clone() (CharSet, error) {
	sc := make([]rune, c1.SetSize())
	copy(sc, c1.setChars)

	return NewCharSet(sc)
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

func (c CharSet) GenRdmBytes(length int) ([]byte, error) {
	if !c.IsAllSingleByte {
		return nil, ErrContainsMultibyteChar
	}

	index := make([]uint8, length)
	if _, err := crand.Read(index); err != nil {
		return nil, err
	}

	randomBytes := make([][]byte, length)
	for i, j := range index {
		randomBytes[i] = c.GetByte(j)
	}

	return bytes.Join(randomBytes, []byte{}), nil
}

func (c CharSet) GenRdmRunes(length int) ([]rune, error) {
	index := make([]uint8, length)

	if c.Rand == nil {
		if _, err := crand.Read(index); err != nil {
			return nil, err
		}
	} else {
		if _, err := c.Rand.Read(index); err != nil {
			return nil, err
		}
	}

	randomRunes := make([]rune, length)
	for i, j := range index {
		randomRunes[i] = c.GetRune(j)
	}
	return randomRunes, nil
}

func (c CharSet) GenRdmString(length int) (string, error) {
	rns, err := c.GenRdmRunes(length)
	return string(rns), err
}
