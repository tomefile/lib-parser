# Tomefile Parser

Library to parse [Tomefile](https://github.com/tomefile) code and output a node tree for use by another program.

<!-- vim-markdown-toc GFM -->

* [Features](#features)
    * [Parsing](#parsing)
    * [Hooks](#hooks)
* [Roadmap](#roadmap)
* [Usage](#usage)

<!-- vim-markdown-toc -->

## Features

### Parsing

Parses the input file into a `*libparser.NodeRoot{}` including string segments.

### Hooks

Allow to run custom `libparser.Hook()` functions on `libparser.Node` before it gets appended to the tree. Returns as soon as an error is encountered. Used to validate, discard or modify nodes.

## Roadmap

Things that need to be done before `v1`:

- [x] Save context to Nodes and allow for partial parsing.
- [ ] Partial parsing
- [ ] Improve runtime string modifier evaluation (currently it just fails silently)
- [ ] Multi-line arguments `(...)`
- [ ] Revisit all of the code and make sure it's well-made (it isn't)
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

parser := libparser.New(file)
parser.Hooks = []libparser.Hook{
    libparser.NoShebangHook,  // remove UNIX shebang, e.g. #!/bin/tome
    libparser.ExcludeHook[*libparser.NodeComment],  // by default, the parser includes ALL file contents, you can discard what's not needed.
}

if derr := parser.Run(); derr != nil {
    derr.Print(os.Stderr)
    os.Exit(1)
}

// [parser.Result] is [*libparser.NodeTree]
return parser.Result
```
