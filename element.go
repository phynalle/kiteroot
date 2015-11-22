package kiteroot

import (
	"bytes"
	"fmt"
	"strings"
)

var selfClosingTagList = []string{
	"area", "base", "br", "col", "command", "embed", "hr", "img", "input",
	"keygen", "link", "meta", "param", "source", "track", "wbr",
}
var selfClosingTagMap map[string]bool

// These are types for Element
const (
	DocumentType ElementType = iota
	TagType
	TextType
)

// LineSeparator is the carrage return.
const LineSeparator = "\n"

// An Attributes stores key/value attribute pairs in the tag.
type Attributes map[string]string

// An ElementType is the type of Element.
type ElementType uint32

// An Element consists of ElementType and Content
type Element struct {
	Type     ElementType
	Content  string
	Attrs    Attributes
	Children []*Element

	selfClosing bool
}

// NewDocument returns a new document type element.
func NewDocument() *Element {
	return &Element{
		Type: DocumentType,
	}
}

// NewTag returns tag typed element for the given name and selfClosing.
func NewTag(name string, selfClosing bool) *Element {
	return &Element{
		Type:        TagType,
		Content:     name,
		Attrs:       make(map[string]string),
		selfClosing: selfClosing,
	}
}

// NewText returns a new text type element for the given content.
func NewText(content string) *Element {
	return &Element{
		Type:    TextType,
		Content: content,
	}
}

// Append appends an element into its children slice in this element.
func (e *Element) Append(elem *Element) {
	if e == elem {
		return
	}
	e.Children = append(e.Children, elem)
}

// SetAttribute puts the given key-value pair into the attribute map in this element.
func (e *Element) SetAttribute(key, value string) {
	e.Attrs[key] = value
}

// Attribute returns the value matching with the key.  if the Attrs doesn't have the key,
// empty string is returned.
func (e *Element) Attribute(key string) string {
	return e.Attrs[key]
}

// FindWithAttrs returns an element containing attrs with the same tag in the subelements.
// if no such element exists, returns nil.
func (e *Element) FindWithAttrs(tagName string, attrs Attributes) *Element {
	return e.findOne(tagName, attrs)
}

// Find returns an element containing attrs with the same tag throughout its subelements.
// if no element is found, this function returns nil.
// attrs is mapped to Attributes by calling MakeAttrs
func (e *Element) Find(tagName string, attrs ...string) *Element {
	return e.findOne(tagName, MakeAttrs(attrs...))
}

// FindAllWithAttrs returns all elements containing attrs with the same tag in the subelements.
func (e *Element) FindAllWithAttrs(tagName string, attrs Attributes) []*Element {
	return e.findAll(tagName, attrs)
}

// FindAll returns all elements containing attrs with the same tag in the subelements.
// attrs is mapped to Attributes by calling MakeAttrs
func (e *Element) FindAll(tagName string, attrs ...string) []*Element {
	return e.findAll(tagName, MakeAttrs(attrs...))
}

// String returns a content according to its type.
func (e *Element) String() string {
	switch e.Type {
	case DocumentType:
		return e.makeDocumentContent()
	case TagType:
		return e.makeTagContent()
	case TextType:
		return e.makeTextContent()
	}
	return ""
}

// Text returns a concatenated text of TextType elements in the children.
func (e *Element) Text() string {
	var contents []string
	for _, child := range e.Children {
		if child.Type == TextType {
			contents = append(contents, child.String())
		}
	}
	return strings.Join(contents, "")
}

// findOne find an element containing the attributes with the same name.
func (e *Element) findOne(name string, attrs Attributes) *Element {
	if e.Type == TextType {
		return nil
	}
	if e.Content == name && e.containsAttrs(attrs) {
		return e
	}
	for _, child := range e.Children {
		founds := child.findOne(name, attrs)
		if founds != nil {
			return founds
		}
	}
	return nil
}

// findAll returns a slice of elements containing the attributes with the same name.
func (e *Element) findAll(name string, attrs Attributes) (tags []*Element) {
	if e.Type == TextType {
		return
	}
	if e.Content == name && e.containsAttrs(attrs) {
		tags = append(tags, e)
	}
	for _, child := range e.Children {
		founds := child.findAll(name, attrs)
		tags = append(tags, founds...)
	}
	return
}

// containsAttrs returns true if the attributes of the element contains
// all of the given attributes.
func (e *Element) containsAttrs(attrs Attributes) bool {
	return containsAttrs(e.Attrs, attrs)
}

// makeDocumentContent returns a concatenated text of its children elements.
func (e *Element) makeDocumentContent() string {
	return e.makeChildrenText()
}

// makeTagContent returns a tag-formatted string. 
// if this element is self-closing tag, its children elements is ignored.
func (e *Element) makeTagContent() string {
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
		buf.WriteString(">")
		buf.WriteRune('\n')
		buf.WriteString(e.makeChildrenText())
		buf.WriteRune('\n')
		buf.WriteString(fmt.Sprintf("</%s>", e.Content))
	}
	return buf.String()
}

// makeTextContent returns Content in this element.
func (e *Element) makeTextContent() string {
	return e.Content
}

// makeChildrenContent concatenates the string of its children, and returns it.
func (e *Element) makeChildrenText() string {
	var contents []string
	for _, child := range e.Children {
		contents = append(contents, child.String())
	}
	return strings.Join(contents, LineSeparator)
}

// MakeAttrs returns Attribute consisting of pairs.
// this function needs the even number of strings to make key/value pairs.
// so, if length of string slice is odd number, the last is dropped.
func MakeAttrs(s ...string) (attrs Attributes) {
	attrs = make(Attributes)
	
	if len(s)%2 == 1 {
		s = s[:len(s)-1]
	}

	for i := 0; i < len(s); i += 2 {
		attrs[s[i]] = s[i+1]
	}
	return
}

// containsAttrs returns true if the given base attributes contains all of identical attributes in attrs.
func containsAttrs(base Attributes, attrs Attributes) bool {
	for key, val := range attrs {
		if v, ok := base[key]; !ok || v != val {
			return false
		}
	}
	return true
}

func init() {
	// Move self-closing-tag list to map because it is more efficient to find a tag.
	selfClosingTagMap = make(map[string]bool)
	for _, tn := range selfClosingTagList {
		selfClosingTagMap[tn] = true
	}
}
