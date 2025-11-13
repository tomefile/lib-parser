package libparser

import "strings"

type Node interface {
	Node() string
}

type NodeTree struct {
	Tomes map[string]Node

	NodeChildren
}

type NodeArgs []Node

func (args NodeArgs) String() string {
	var builder strings.Builder

	for _, arg := range args {
		builder.WriteString(" " + arg.Node())
	}

	return builder.String()
}

type NodeChildren []Node

func (children NodeChildren) String() string {
	if len(children) == 0 {
		return " {}"
	}

	var builder strings.Builder
	builder.WriteString(" {\n")

	for _, arg := range children {
		builder.WriteString("--- " + arg.Node() + "\n")
	}

	builder.WriteString("}")
	return builder.String()
}
