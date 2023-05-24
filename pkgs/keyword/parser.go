package keyword

type KeywordReader struct{}

type Token struct {
	Op   Operater
	Text string
}

func (r KeywordReader) Lex(str string) ([]Token, error)

func (r KeywordReader) IsOperator(rn rune) bool {
	if rn == '+' || rn == '-' {
		return true
	}
	return false
}

func (r KeywordReader) IsQutation(rn rune) bool {
	return rn == '"'
}

func (r KeywordReader) Parse()
