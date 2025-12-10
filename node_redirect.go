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
	node.printIfNotNil(&builder, node.Stdin, " <")
	node.printIfNotNil(&builder, node.Stdout, " >")
	node.printIfNotNil(&builder, node.Stderr, " >>")
	return node.Source.String() + builder.String()
}

func (node *NodeRedirect) printIfNotNil(b *strings.Builder, target *NodeString, prefix string) {
	if target != nil {
		b.WriteString(prefix + target.String())
	}
}
