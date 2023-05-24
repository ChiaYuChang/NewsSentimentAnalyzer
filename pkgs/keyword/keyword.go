package keyword

import "fmt"

type Keyword interface {
	ToString() string
}

type Operater string

const (
	AND Operater = "AND"
	OR  Operater = "OR"
	NOT Operater = "NOT"
)

var EmptyElement = Element("")

type KeywordsRelation struct {
	LHS, RHS Keyword
	Op       Operater
}

func NewKeywordsRelation(lhs, rhs Keyword, op Operater) KeywordsRelation {
	if op == NOT {
		return KeywordsRelation{LHS: EmptyElement, RHS: rhs, Op: op}
	}
	return KeywordsRelation{LHS: lhs, RHS: rhs, Op: op}
}

func (kr KeywordsRelation) ToString() string {
	if kr.Op == NOT {
		return fmt.Sprintf("%s (%s)", NOT, kr.RHS.ToString())
	}
	return fmt.Sprintf("(%s %s %s)", kr.LHS, kr.Op, kr.RHS)
}
