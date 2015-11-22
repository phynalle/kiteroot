package kiteroot

import (
	"reflect"
	"testing"
)

func TestMakeAttrs(t *testing.T) {
	testcase := [][]string{
		[]string{"key", "value", "id", "post"},
		[]string{"src", "somewhere"},
		[]string{"href", "#", "type", "text/css", "rel", "stylesheet", "dropped"},
		[]string{"hash", "1234", "class", "secret"},
	}

	attrs := []Attributes{
		Attributes{
			"key": "value",
			"id":  "post",
		},
		Attributes{
			"src": "somewhere",
		},
		Attributes{
			"href": "#",
			"rel":  "stylesheet",
			"type": "text/css",
		},
		Attributes{
			"hash":  "1234",
			"class": "secret",
		},
	}

	for i, tc := range testcase {
		if !reflect.DeepEqual(MakeAttrs(tc...), attrs[i]) {
			t.Fail()
		}
	}
}

func TestAttribute(t *testing.T) {
	baseAttrs := Attributes{
		"id":        "site-head",
		"class":     "header",
		"title":     "attribute",
		"condition": "yes",
	}

	if !containsAttrs(baseAttrs, nil) {
		t.Errorf("Contains nil Attributes should be true")
	}
	testcase := [][]string{
		[]string{"id", "site-head"},
		[]string{"class", "content"},
		[]string{"title", "attribute", "condition", "yes"},
		[]string{"id", "site-head", "class", "header", "title", "attribute"},
	}
	results := []bool{true, false, true, true}

	for i, tc := range testcase {
		attrs := MakeAttrs(tc...)
		if containsAttrs(baseAttrs, attrs) != results[i] {
			t.Fail()
		}
	}
}

func TestTagString(t *testing.T) {
	var text Element
	text.Type = TextType
	text.Content = "text"

	var tag Element
	tag.Type = TagType
	tag.Content = "name"
	tag.Attrs = MakeAttrs("class", "attribute", "attr", "yes")
	tag.Append(&text)

	var br Element
	br.Type = TagType
	br.Content = "br"
	br.Attrs = MakeAttrs("class", "correct")
	br.selfClosing = true

	tag.Append(&br)
	found := tag.FindAll("br", "class", "correct")
	if len(found) != 1 {
		t.Fail()
	}
	elem := tag.FindWithAttrs("br", MakeAttrs("class", "wrong"))
	if elem != nil {
		t.Fail()
	}

}
