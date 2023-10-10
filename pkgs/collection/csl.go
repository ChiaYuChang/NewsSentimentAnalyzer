package collection

import (
	"regexp"
	"strconv"
	"strings"
)

// Comma Separated List
type CSL []string

// Check if string is unicode escaped
func IsUnicodeEscaped(s string) bool {
	return strings.Contains(s, "\\u")
}

// Create new comma separated list from string (eg. "a,b,c")
func NewCSL(s string) CSL {
	ss := strings.Split(s, ",")
	csl := make([]string, 0, len(ss))
	for _, s := range ss {
		if s = strings.TrimSpace(s); s != "" {
			if IsUnicodeEscaped(s) {
				if uqs, err := strconv.Unquote(`"` + s + `"`); err == nil {
					s = uqs
				}
			}
			csl = append(csl, s)
		}
	}
	return csl
}

func (csl CSL) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strings.Join(csl, ","))), nil
}

func (csl *CSL) UnmarshalJSON(b []byte) error {
	s := strings.ReplaceAll(string(b), "\n", "")
	if s[0] == '[' && s[len(s)-1] == ']' {
		var ss []string
		re := regexp.MustCompile("\"([^,]+?)\"")
		for _, t := range re.FindAllStringSubmatch(s, -1)[1:] {
			ss = append(ss, string(t[1]))
		}
		(*csl) = CSL(ss)
	} else {
		s = strings.Trim(s, "\"")
		(*csl) = NewCSL(s)
	}
	return nil
}
