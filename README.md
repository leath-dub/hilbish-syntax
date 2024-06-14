# Hilbish Syntax

This is a native hilbish module that provides syntax highlighting.

## Setup

You need to first build the native module, and for that you
will need golang:

```sh
go build -buildmode=plugin
```

This will create a `hilbish-syntax.so` file in the current directory
(Assuming your on a unix-like OS, honestly don't care if your Windows: you get enough support from the big players)

Then to enable it in your config:

```lua
hilbish.highlighter = hilbish.module.load(os.getenv 'HOME' .. '/Repos/hilbish-syntax/hilbish-syntax.so')
```

where `os.getenv 'HOME' .. /Repos/hilbish-syntax/hibish-syntax.so` is the path to the `.so` file (relative to `$HOME` dir)

## For the hackers !

This project is built using tree-sitter. You can define your own highlight rules in the `highlights.scm` file which is built
into the `.so` produced. You can consult the tree-sitter docs if you want to create custom highlights, but importantly any of
the `@<name>` variables can be assigned a color in the source code, here is a snippet from `syntax.go`:

```go
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
```

This is just a map of those variables to an ansi code which is then inserted at the start of the nodes captured by the variable
(a `\033[0m` is matched at the end of the nodes captured)

## NOTE

This is very early stages and currently only supports posix-like command-line highlighting.
Additionally this has some latency issues, more profiling needs to be done to see what can
be done to reduce the latency however there seems to be a fixed burden of calling between
the native module and hilbish.
