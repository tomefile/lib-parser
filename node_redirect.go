package libparser

import "strings"

type NodeRedirect struct {
	Source Node
	Stdin  *NodeString
	Stdout *NodeString
	Stderr *NodeString
	NodeContext
}

func (node *NodeRedirect) Context() NodeContext {
	return node.NodeContext
}

func (node *NodeRedirect) String() string {
	var builder strings.Builder
	if node.Stdin != nil {
		builder.WriteString(" <" + node.Stdin.Segments.String())
	}
	if node.Stdout != nil {
		builder.WriteString(" >" + node.Stdout.Segments.String())
	}
	if node.Stderr != nil {
		builder.WriteString(" >>" + node.Stderr.Segments.String())
	}
	return node.Source.String() + builder.String()
}
