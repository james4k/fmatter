package fmatter

import (
	"bytes"
	"fmt"
	"testing"
)

func ExampleRead() {
	data := []byte(`---
title: Some Title
---
content`)

	var frontmatter struct{ Title string }
	content, _ := Read(data, &frontmatter)
	fmt.Println(content)
}

// TODO: test a few expected errors

type testItem struct {
	data            []byte
	expectedContent []byte
}

var testItems = []testItem{
	{[]byte(`---
frontmatter: simple
---
content`),
		[]byte(`content`)},
	{[]byte(`  content`),
		[]byte(`  content`)},
	{[]byte(`---
content`),
		[]byte(`---
content`)},
}

func TestItems(t *testing.T) {
	frontmatter := make(map[string]interface{})

	for _, item := range testItems {
		content, err := Read(item.data, frontmatter)
		if err != nil {
			t.Fatal(err)
		}

		if bytes.Compare(content, item.expectedContent) != 0 {
			t.Fatalf("unexpected content:\n%v\nvs.\n%v",
				string(content), string(item.expectedContent))
		}
	}
}
