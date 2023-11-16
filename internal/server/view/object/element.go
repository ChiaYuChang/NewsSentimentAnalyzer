package object

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"strings"
)

var ErrUnknownAttrType = errors.New("unknown attr type")

type HTMLElementList struct {
	Tag string         `json:"tag"`
	Ele []*HTMLElement `json:"elements"`
}

func NewHTMLElementList(tag string) *HTMLElementList {
	return &HTMLElementList{
		Tag: tag,
		Ele: make([]*HTMLElement, 0),
	}
}

func (el *HTMLElementList) NewHTMLElement(attrs ...HTMLAttr) *HTMLElement {
	e := NewHTMLElement(el.Tag, attrs...)
	(*el).Ele = append((*el).Ele, e)
	return e
}

func (el HTMLElementList) Element() []*HTMLElement {
	return el.Ele
}

func (el1 *HTMLElementList) Copy() *HTMLElementList {
	el2 := NewHTMLElementList(el1.Tag)
	el2.Ele = make([]*HTMLElement, len(el1.Ele))
	for i, e := range el1.Ele {
		el2.Ele[i] = e.Copy()
	}
	return el2
}

type HTMLElement struct {
	Attr          []HTMLAttr    `json:"attr"`
	IsSelfClosing bool          `json:"id_self_closing"`
	Content       template.HTML `json:"content"`
	Tag           string        `json:"tag"`
}

func (e *HTMLElement) UnmarshalJSON(data []byte) error {
	type InnerHTMLElement HTMLElement

	el := struct {
		*InnerHTMLElement
		Attr []any `json:"attr"`
	}{}

	err := json.Unmarshal(data, &el)
	if err != nil {
		return err
	}

	for _, a := range el.Attr {
		switch v := a.(type) {
		case map[string]any:
			el.InnerHTMLElement.Attr = append(
				el.InnerHTMLElement.Attr, HTMLAttrPair{
					Tag: v["tag"].(string),
					Val: v["val"].(string),
				})
		case string:
			el.InnerHTMLElement.Attr = append(
				el.InnerHTMLElement.Attr, HTMLAttrVal(v))
		default:
			return ErrUnknownAttrType
		}
	}

	(*e) = HTMLElement(*el.InnerHTMLElement)
	return nil
}

type HTMLAttr interface {
	ToHTMLAttr() template.HTMLAttr
	String() string
	Copy() HTMLAttr
}

func NewHTMLElement(tag string, attrs ...HTMLAttr) *HTMLElement {
	e := HTMLElement{Attr: attrs, Tag: tag}
	return &e
}

func (e *HTMLElement) ToSelfClosingElement() *HTMLElement {
	e.Content = ""
	e.IsSelfClosing = true
	return e
}

func (e *HTMLElement) ToOpeningElement(content template.HTML) *HTMLElement {
	e.Content = content
	e.IsSelfClosing = false
	return e
}

func (e *HTMLElement) AddPair(key string, val string) *HTMLElement {
	e.Attr = append(e.Attr, HTMLAttrPair{Tag: key, Val: val})
	return e
}

func (e *HTMLElement) AddVal(val string) *HTMLElement {
	e.Attr = append(e.Attr, HTMLAttrVal(val))
	return e
}

func (e HTMLElement) ToHTMLAttr() template.HTMLAttr {
	return template.HTMLAttr(e.String())
}

func (e HTMLElement) ToHTML() template.HTML {
	if e.IsSelfClosing {
		return template.HTML(fmt.Sprintf("<%s %s />", e.Tag, e.String()))
	}
	return template.HTML(fmt.Sprintf("<%s %s> %s </%s>", e.Tag, e.String(), e.Content, e.Tag))
}

func (e1 *HTMLElement) Copy() *HTMLElement {
	e2 := &HTMLElement{}
	e2.Attr = make([]HTMLAttr, len(e1.Attr))
	for i, a := range e1.Attr {
		e2.Attr[i] = a.Copy()
	}
	e2.Content = e1.Content
	e2.IsSelfClosing = e1.IsSelfClosing
	e2.Tag = e1.Tag
	return e2
}

func (e HTMLElement) String() string {
	attrs := make([]string, len(e.Attr))
	for i := range attrs {
		attrs[i] = e.Attr[i].String()
	}
	return strings.Join(attrs, " ")
}

type HTMLAttrPair struct {
	Tag string `json:"tag"`
	Val string `json:"val"`
}

func (p HTMLAttrPair) ToHTMLAttr() template.HTMLAttr {
	return template.HTMLAttr(p.String())
}

func (p HTMLAttrPair) String() string {
	return fmt.Sprintf("%s=\"%s\"", p.Tag, p.Val)
}

func (p HTMLAttrPair) Copy() HTMLAttr {
	return HTMLAttrPair{Tag: p.Tag, Val: p.Val}
}

type HTMLAttrVal string

func (v HTMLAttrVal) ToHTMLAttr() template.HTMLAttr {
	return template.HTMLAttr(v.String())
}

func (v HTMLAttrVal) String() string {
	return string(v)
}

func (v HTMLAttrVal) Copy() HTMLAttr {
	return HTMLAttrVal(v)
}
