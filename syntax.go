package main

import (
	"context"
	_ "embed"
	"strings"

	rt "github.com/arnodel/golua/runtime"
	ts "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/bash"
)

//go:embed highlights.scm
var highlights []byte

//go:embed test.sh
var test []byte

func Loader(rtm *rt.Runtime) rt.Value {
	hl, err := NewHighlighter(bash.GetLanguage())
	if err != nil {
		panic(err)
	}
	fn := func(tt *rt.Thread, cc *rt.GoCont) (rt.Cont, error) {
		if err := cc.Check1Arg(); err != nil {
			return nil, err
		}

		ln, err := cc.StringArg(0)
		if err != nil {
			return nil, err
		}

		ln, err = hl.Highlight(ln)
		if err != nil {
			return nil, err
		}

		return cc.PushingNext1(tt.Runtime, rt.StringValue(string(ln))), nil
	}
	r := rtm.SetEnvGoFunc(rtm.GlobalEnv(), "hilbish.highlighter", fn, 1, false)
	return rt.FunctionValue(r)
}

type Highlighter struct {
	lang    *ts.Language
	oldTree *ts.Tree
	parser  *ts.Parser
	query   *ts.Query
}

func NewHighlighter(lang *ts.Language) (*Highlighter, error) {
	q, err := ts.NewQuery(highlights, lang)
	if err != nil {
		return nil, err
	}
	parser := ts.NewParser()
	parser.SetLanguage(lang)
	return &Highlighter{lang: lang, oldTree: nil, parser: parser, query: q}, nil
}

var THEME = map[string]string{
	"string":   "\033[38;5;3m",
	"function": "\033[38;5;13m",
	"property": "\033[38;5;15m",
	"keyword":  "\033[38;5;9m",
	"number":   "\033[38;5;5m",
	"embedded": "\033[38;5;3m",
	"operator": "\033[38;5;11m",
	"constant": "\033[38;5;1m",
}

func (h *Highlighter) Highlight(text string) (string, error) {
	tree, err := h.parser.ParseCtx(context.Background(), h.oldTree, []byte(text))
	if err != nil {
		return "", err
	}

	root := tree.RootNode()

	qc := ts.NewQueryCursor()
	qc.Exec(h.query, root)

	type Highlight struct {
		isSome  bool
		code    string
		endByte uint32
	}

	highlights := make([]Highlight, len(text))

	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		m = qc.FilterPredicates(m, []byte(text))

		for _, c := range m.Captures {
			deco := h.query.CaptureNameForId(c.Index)
			code, ok := THEME[deco]
			if !ok {
				continue
			}

			sb, eb := c.Node.StartByte(), c.Node.EndByte()
			highlights[sb] = Highlight{isSome: true, code: code, endByte: eb}
		}
	}

	var bld strings.Builder
	var ofs uint32

	for ofs = 0; ofs < uint32(len(text)); ofs++ {
		hl := highlights[ofs]
		if !hl.isSome {
			bld.WriteByte(text[ofs])
			continue
		}

		bld.WriteString(hl.code)
		bld.WriteString(string(text[ofs:hl.endByte]))
		bld.WriteString("\033[0m")

		ofs = hl.endByte - 1 // skip the rest of the highlight
	}

	return bld.String(), nil
}
