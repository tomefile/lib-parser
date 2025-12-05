package libparser

import "strings"

type NodeRedirectOut struct {
	Source    Node
	Filenames []*NodeString
	NodeContext
}

func (node *NodeRedirectOut) Context() NodeContext {
	return node.NodeContext
}

func (node *NodeRedirectOut) String() string {
	var builder strings.Builder
	for _, filename := range node.Filenames {
		builder.WriteString(" > " + filename.String())
	}
	return node.Source.String() + builder.String()
}

// ————————————————————————————————

type NodeHereString struct {
	Source *NodeString
	Dest   Node
	NodeContext
}

func (node *NodeHereString) Context() NodeContext {
	return node.NodeContext
}

func (node *NodeHereString) String() string {
	return node.Dest.String() + " << " + node.Source.String()
}
