# Tomefile Parser

Library to parse [Tomefile](https://github.com/tomefile) code and output a node tree for use by another program.

<!-- vim-markdown-toc GFM -->

* [Features](#features)
* [Usage](#usage)

<!-- vim-markdown-toc -->

## Features

- **Parsing** `libparser.Parse(...)` — Parses the input UTF-8 stream into a node tree
- **Formatting** `libparser.Format(...)` — Formats the input UTF-8 stream returning a slice of segments. Used to substitute environmental, local, and such variables.
- **Post-Processing** `func(Node) (Node, *DetailedError)` — Inject functions into `libparser.Parse(...)` to be applied to a node before it gets appended to the tree. Returns as soon as an error is encountered. Used to validate, discard or modify nodes.

## Usage

Parsing a file:

```go
parser := libparser.New(
    "example.tome",
    bufio.NewReader(file),
    libparser.PostNoShebang,  // remove UNIX shebang, e.g. #!/bin/tome
    libparser.PostExclude[*libparser.DirectiveNode], // let's say we want to exclude a specific node type
)
tree, err := parser.Parse()
if err != nil {
    err.BeautyPrint(os.Stderr)
    os.Exit(1)
}

// [tree] is [*libparser.NodeTree]
```

Formatting a variable (i.e. `$name ${name:mod} etc.`)

```go
formatter := libparser.NewStringFormatter(
    bufio.NewReader(strings.NewReader(
        "this is an example $string with ${string:trim_suffix 123}",
    )),
)

segments, err := formatter.Format()
if err != nil {
    err.BeautyPrint(os.Stderr)
    os.Exit(1)
}

// [segments] is [[]libparser.Segment]
```
