package htmltable

import (
	"encoding/xml"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type ColTypeHandlerRepo map[DType]DTypeHandler

func Less[T int | string](x, y T, order bool) bool {
	if x != y {
		if order {
			return x >= y
		} else {
			return x < y
		}
	}
	return false
}

var repo ColTypeHandlerRepo = ColTypeHandlerRepo{
	DTypeString: func(x, y string, order bool) bool {
		return Less(x, y, order)
	},
	DTypeInt: func(x, y string, order bool) bool {
		ix, _ := strconv.Atoi(x)
		iy, _ := strconv.Atoi(y)
		return Less[int](ix, iy, order)
	},
}

type DType int8

const (
	DTypeString DType = iota
	DTypeInt
)
const DTypeDefault = DTypeString

type DTypeHandler func(x, y string, order bool) bool

func (dt DType) String() string {
	var s string
	switch dt {
	case DTypeString:
		s = "string"
	case DTypeInt:
		s = "int"
	default:
		s = "unknown"
	}
	return s
}

func ParseDType(s string) DType {
	var dt DType
	switch s {
	case "string":
		dt = DTypeString
	case "int":
		dt = DTypeInt
	default:
		dt = DTypeDefault
	}
	return dt
}

func (s *DType) UnmarshalText(text []byte) error {
	*s = ParseDType(string(text))
	return nil
}

func (s DType) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

const (
	ASC = false
	DES = true
)

var ErrIndexAndOrderArrayLengthNotMatch = errors.New("index and order array length not match")

type HTMLTable struct {
	cache   []bool        `xml:"-"`
	Class   string        `xml:"class,attr"`
	XMLName xml.Name      `xml:"table"`
	Head    HTMLTableHead `xml:"thead"`
	Body    HTMLTableBody `xml:"tbody"`
}

type HTMLTableHead struct {
	Row []struct {
		Header []struct {
			Text    string `xml:",innerxml"`
			DType   DType  `xml:"dtype,attr"`
			OnClick string `xml:"onclick,attr,omitempty"`
		} `xml:"th"`
	} `xml:"tr"`
	XMLName xml.Name `xml:"thead"`
}

type HTMLTableBody struct {
	Row []struct {
		Class string `xml:"class,attr,omitempty"`
		Data  []struct {
			Text string `xml:",innerxml"`
		} `xml:"td"`
	} `xml:"tr"`
	XMLName xml.Name `xml:"tbody"`
}

func (tb HTMLTable) String() string {
	sb := strings.Builder{}
	sb.WriteString("HTML Table:\n")
	sb.WriteString("\t- Header:\n")
	for _, header := range tb.Head.Row[0].Header {
		sb.WriteString(fmt.Sprintf("\t\t- %s (%s)\n", header.Text, header.DType.String()))
	}

	sb.WriteString("\t- Body:\n")
	for _, row := range tb.Body.Row {
		data := make([]string, tb.NRow())
		for i, d := range row.Data {
			data[i] = d.Text
		}
		sb.WriteString("\t\t- " + strings.Join(data, ", ") + "\n")
	}
	return sb.String()
}

func (tb HTMLTable) NCol() int {
	return len(tb.Head.Row[0].Header)
}

func (tb HTMLTable) NRow() int {
	return len(tb.Body.Row)
}

func (tb HTMLTable) DType(i int) DType {
	if i >= tb.NCol() {
		return DTypeDefault
	}
	return tb.Head.Row[0].Header[i].DType
}

func (tb *HTMLTable) newCache() {
	if len(tb.cache) == 0 {
		(tb.cache) = make([]bool, tb.NCol())
	}
}

func (tb *HTMLTable) SortBy(c int, o bool) *HTMLTable {
	if c > tb.NCol() {
		return tb
	}
	tb.newCache()
	tb.cache[c] = o
	tbsb := HTMLTableSortBy{tb, c, tb.cache[c]}
	sort.Sort(tbsb)
	return tbsb.HTMLTable
}

func (tb *HTMLTable) SortByWithMemory(c int) *HTMLTable {
	if c > tb.NCol() {
		return tb
	}

	tb.newCache()
	tb.cache[c] = !tb.cache[c]
	tbsb := HTMLTableSortBy{tb, c, tb.cache[c]}
	sort.Sort(tbsb)
	return tbsb.HTMLTable
}

type HTMLTableSortBy struct {
	*HTMLTable
	Index int
	Order bool
}

func (tb HTMLTableSortBy) Len() int {
	return tb.NRow()
}

func (tb HTMLTableSortBy) Swap(i, j int) {
	tb.Body.Row[i], tb.Body.Row[j] = tb.Body.Row[j], tb.Body.Row[i]
}

func (tb HTMLTableSortBy) Less(i, j int) bool {
	var handler DTypeHandler
	var ok bool

	if handler, ok = repo[tb.DType(tb.Index)]; !ok {
		handler = repo[DTypeString]
	}

	return handler(
		tb.Body.Row[i].Data[tb.Index].Text,
		tb.Body.Row[j].Data[tb.Index].Text,
		tb.Order,
	)
}
