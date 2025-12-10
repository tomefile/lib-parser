# Tomefile Parser

Library to parse [Tomefile](https://github.com/tomefile) code and output a node tree for use by another program.

<!-- vim-markdown-toc GFM -->

* [Features](#features)
* [Roadmap](#roadmap)
* [Usage](#usage)

<!-- vim-markdown-toc -->

## Features

- **Parsing** `libparser.Parse(...)` — Parses the input UTF-8 stream into a node tree
- **Formatting** `libparser.Format(...)` — Formats the input UTF-8 stream returning a slice of segments. Used to substitute environmental, local, and such variables.
- **Post-Processing** `func(Node) (Node, *liberrors.DetailedError)` — Inject functions into `libparser.Parse(...)` to be applied to a node before it gets appended to the tree. Returns as soon as an error is encountered. Used to validate, discard or modify nodes.

## Roadmap

Things that need to be done before `v1`:

- [x] Save context to Nodes and allow for partial parsing.
- [ ] Partial parsing
- [ ] Multi-line arguments `(...)`
- [ ] Better test coverage (Add more edge-cases where formatting isn't perfect)

## Usage

Parsing a file:

```go
// Close all files when all parsers have finished.
defer libparser.CloseAll()

file, err := libparser.OpenFile("example.tome")
if err != nil {
    panic(err)
}

parser := libparser.New(file).
    With(libparser.PostNoShebang).  // remove UNIX shebang, e.g. #!/bin/tome
    With(libparser.PostExclude[*libparser.CommentNode])  // let's say we want to exclude a specific node type

tree, derr := parser.Parse()
if derr != nil {
    derr.Print(os.Stderr)
    os.Exit(1)
}

// [tree] is [*libparser.NodeTree]
```

Formatting a variable (i.e. `$name ${name:mod} etc.`)

```go
formatter := libparser.NewStringFormatter("this is an example $string with ${string:trim_suffix 123}",)

segments, derr := formatter.Format()
if derr != nil {
    derr.Print(os.Stderr)
    os.Exit(1)
}

// [segments] is [[]libparser.Segment]
```

Add RedirectSource
