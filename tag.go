package kiteroot

import (
	"bytes"
	"fmt"
	_ "golang.org/x/net/html"
	"strings"
)

var selfClosingTagList = []string{
	"area", "base", "br", "col", "command", "embed", "hr", "img", "input",
	"keygen", "link", "meta", "param", "source", "track", "wbr",
}
var selfClosingTagMap map[string]bool

const (
	DocumentType ElementType = iota
	TagType
	TextType
)

const (
	LineSeparator = "\n"
)

type ElementType uint32

type Element struct {
	Type     ElementType
	Content  string
	Attrs    Attributes
	Children []*Element

	selfClosing bool
}

func NewDocument() *Element {
	return &Element{
		Type: DocumentType,
	}
}

func NewTag(name string, sc bool) *Element {
	return &Element{
		Type:        TagType,
		Content:     name,
		Attrs:       make(map[string]string),
		selfClosing: sc,
	}
}

func NewText(content string) *Element {
	return &Element{
		Type:    TextType,
		Content: content,
	}
}

func (e *Element) Append(elem *Element) {
	if e == elem {
		return
	}
	e.Children = append(e.Children, elem)
}

func (e *Element) SetAttribute(key, value string) {
	e.Attrs[key] = value
}

func (e *Element) Attribute(key string) (string, bool) {
	val, ok := e.Attrs[key]
	return val, ok
}

func (e *Element) Find(name string, attrs Attributes) (tags []*Element) {
	if e.Type == TextType {
		return
	}

	if e.Content == name && e.containsAttrs(attrs) {
		tags = append(tags, e)
	}

	for _, child := range e.Children {
		founds := child.Find(name, attrs)
		tags = append(tags, founds...)
	}
	return
}

func (e *Element) String() string {
	switch e.Type {
	case DocumentType:
		return e.toDocumentString()
	case TagType:
		return e.toTagString()
	case TextType:
		return e.toTextString()
	}
	return ""
}

func (e *Element) containsAttrs(attrs Attributes) bool {
	return containsAttrs(e.Attrs, attrs)
}

func (e *Element) toDocumentString() string {
	var contents []string
	for _, child := range e.Children {
		contents = append(contents, child.String())
	}
	return strings.Join(contents, LineSeparator)
}

func (e *Element) toTagString() string {
	var attrs []string
	for k, v := range e.Attrs {
		attrs = append(attrs, fmt.Sprintf("%s=\"%s\"", k, v))
	}

	var buf bytes.Buffer
	buf.WriteRune('<')
	buf.WriteString(e.Content)

	if len(attrs) > 0 {
		buf.WriteRune(' ')
		buf.WriteString(strings.Join(attrs, " "))
	}

	if e.selfClosing {
		buf.WriteString(" />")
	} else {
		var contents []string
		for _, child := range e.Children {
			contents = append(contents, child.String())
		}

		buf.WriteString(">")
		buf.WriteRune('\n')
		buf.WriteString(strings.Join(contents, LineSeparator))
		buf.WriteRune('\n')
		buf.WriteString(fmt.Sprintf("</%s>", e.Content))
	}
	return buf.String()
}

func (e *Element) toTextString() string {
	return e.Content
}

type Attributes map[string]string

func containsAttrs(base Attributes, attrs Attributes) bool {
	for key, val := range attrs {
		if v, ok := base[key]; ok && v != val {
			return false
		}
	}
	return true
}

func MakeAttrs(s ...string) (attrs Attributes) {
	attrs = make(map[string]string)

	if len(s)%2 == 1 {
		s = s[:len(s)-1]
	}

	for i := 0; i < len(s); i += 2 {
		attrs[s[i]] = s[i+1]
	}
	return
}

func init() {
	selfClosingTagMap = make(map[string]bool)
	for _, tn := range selfClosingTagList {
		selfClosingTagMap[tn] = true
	}
}
