package object

import (
	"fmt"
	"html/template"
	"strings"
)

type HTMLElementList struct {
	tag string
	e   []*HTMLElement
}

func NewHTMLElementList(tag string) *HTMLElementList {
	return &HTMLElementList{
		tag: tag,
		e:   make([]*HTMLElement, 0),
	}
}

func (el *HTMLElementList) NewHTMLElement(attrs ...HTMLAttr) *HTMLElement {
	e := NewHTMLElement(el.tag, attrs...)
	(*el).e = append((*el).e, e)
	return e
}

func (el HTMLElementList) Element() []*HTMLElement {
	return el.e
}

func (el1 *HTMLElementList) Copy() *HTMLElementList {
	el2 := NewHTMLElementList(el1.tag)
	el2.e = make([]*HTMLElement, len(el1.e))
	for i, e := range el1.e {
		el2.e[i] = e.Copy()
	}
	return el2
}

type HTMLElement struct {
	Attr          []HTMLAttr
	IsSelfClosing bool
	Content       template.HTML
	Tag           string
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
	Tag, Val string
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
