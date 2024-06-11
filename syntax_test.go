package main

import (
	"github.com/smacker/go-tree-sitter/bash"
	"testing"
)

var hl *Highlighter

func init() {
	var err error
	hl, err = NewHighlighter(bash.GetLanguage())
	if err != nil {
		panic(err)
	}
}

func BenchmarkNewHighlighter(b *testing.B) {
	_, err := NewHighlighter(bash.GetLanguage())
	if err != nil {
		panic(err)
	}
}

func BenchmarkHighlight(b *testing.B) {
	for i := 0; i < 1; i++ {
		_, err := hl.Highlight(string(test))
		if err != nil {
			panic(err)
		}
	}
}
