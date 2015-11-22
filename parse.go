package kiteroot

import (
	"errors"
	"io"
	"strings"

	"golang.org/x/net/html"
)

var (
	ErrInvalidPair = errors.New("open/close tag mismatched")
)

// Parse returns element tree for the HTML from the given Reader.
func Parse(r io.Reader) (*Element, error) {
	var st Stack
	z := html.NewTokenizer(r)
	doc := NewDocument()
	st.Push(doc)

ParseIterator:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				break ParseIterator
			}
			return nil, z.Err()

		case html.StartTagToken:
			fallthrough

		case html.SelfClosingTagToken:
			t := z.Token()
			sc := selfClosingTagMap[t.Data] || tt == html.SelfClosingTagToken
			tag := NewTag(t.Data, sc)

			for _, attr := range t.Attr {
				tag.SetAttribute(attr.Key, attr.Val)
			}

			cur := st.Top()
			if cur == nil {
				return nil, ErrInvalidPair
			}

			cur.Append(tag)
			if !sc {
				st.Push(tag)
			}

		case html.EndTagToken:
			t := z.Token()
			if selfClosingTagMap[t.Data] {
				continue
			}
			if cur := st.Pop(); cur == nil || cur.Content != t.Data {
				return nil, ErrInvalidPair
			}

		case html.TextToken:
			cur := st.Top()
			if cur == nil {
				return nil, ErrInvalidPair
			}
			s := string(z.Text())
			if s == "\n" {
				continue
			}
			text := NewText(strings.TrimSpace(s))
			cur.Append(text)
		}
	}
	if st.Len() != 1 {
		return nil, ErrInvalidPair
	}
	return doc, nil
}
