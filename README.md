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
- **Post-Processing** `func(Node) (Node, *DetailedError)` — Inject functions into `libparser.Parse(...)` to be applied to a node before it gets appended to the tree. Returns as soon as an error is encountered. Used to validate, discard or modify nodes.

## Roadmap

- [ ] Macros `example!`.
- [ ] Support for `;` to separate statements.
- [ ] Support `&&` and `||` in commands.
- [ ] Pipes `|`.
- [ ] Redirects `>` & `<`.

## Usage

Parsing a file:

```go
parser, err := libparser.OpenNew(
    "example.tome",
    libparser.PostNoShebang,  // remove UNIX shebang, e.g. #!/bin/tome
    libparser.PostExclude[*libparser.DirectiveNode], // let's say we want to exclude a specific node type
)
if err != nil {
    panic(err) // File couldn't be opened for reading
}
defer parser.Close()

tree, detailed_err := parser.Parse()
if detailed_err != nil {
    detailed_err.BeautyPrint(os.Stderr)
    os.Exit(1)
}

// [tree] is [*libparser.NodeTree]
```

Formatting a variable (i.e. `$name ${name:mod} etc.`)

```go
formatter := libparser.NewStringFormatter("this is an example $string with ${string:trim_suffix 123}",)

segments, detailed_err := formatter.Format()
if detailed_err != nil {
    detailed_err.BeautyPrint(os.Stderr)
    os.Exit(1)
}

// [segments] is [[]libparser.Segment]
```
