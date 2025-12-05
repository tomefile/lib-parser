package libparser

import (
	"fmt"
	"strings"
)

type Node interface {
	Context() NodeContext
	String() string
}

// ————————————————————————————————

type NodeRoot struct {
	Tomes map[string]Node
	NodeContext
	NodeChildren
}

func (node *NodeRoot) Context() NodeContext {
	return node.NodeContext
}

func (node *NodeRoot) String() string {
	return node.NodeChildren.String()
}

// ————————————————————————————————

type NodeContext struct {
	OffsetStart, OffsetEnd uint
}

func (context NodeContext) String() string {
	return fmt.Sprintf("[%d-%d]", context.OffsetStart, context.OffsetEnd)
}

// ————————————————————————————————

type NodeChildren []Node

func (children NodeChildren) String() string {
	if len(children) == 0 {
		return " {}"
	}

	var builder strings.Builder
	builder.WriteString(" {\n")

	for _, arg := range children {
		builder.WriteString("--- " + arg.String() + "\n")
	}

	builder.WriteString("}")
	return builder.String()
}

// ————————————————————————————————

type NodeArgs []Node

func (args NodeArgs) String() string {
	var builder strings.Builder

	for _, arg := range args {
		builder.WriteString(" " + arg.String())
	}

	return builder.String()
}
