package keyword

import (
	"fmt"
	"strings"
)

type Element string

func (e Element) ToString() string {
	return string(e)
}

type KeywordSet map[string]struct{}

func NewKeywordSet(kws ...Keyword) KeywordSet {
	ks := make(map[string]struct{})
	for _, kw := range kws {
		ks[kw.ToString()] = struct{}{}
	}
	return ks
}

func (ks KeywordSet) Append(kws ...Keyword) (ok bool) {
	for _, kw := range kws {
		kwStr := kw.ToString()
		if _, ok = ks[kwStr]; !ok {
			ks[kwStr] = struct{}{}
		}
	}
	return ok
}

func (ks KeywordSet) ToString(operator string) string {
	kwStrs := make([]string, 0, len(ks))
	for kw := range ks {
		kwStrs = append(kwStrs, kw)
	}
	return fmt.Sprintf("(%s)", strings.Join(kwStrs, operator))
}

func (ks KeywordSet) Intersection() Element {
	return Element(ks.ToString(string(AND)))
}

func (ks KeywordSet) Union() Element {
	return Element(ks.ToString(string(OR)))
}
