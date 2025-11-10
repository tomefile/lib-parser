# Tomefile Parser

Library to parse [Tomefile](https://github.com/tomefile) code and output a node tree for use by another program.

## Usage

Parsing a file:

```go
parser := libparser.New("example.tome", bufio.NewReader(file))
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
