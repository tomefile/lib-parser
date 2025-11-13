package libparser

import "fmt"

type DirectiveNode struct {
	Name string

	NodeArgs
	NodeChildren
}

func (node *DirectiveNode) Node() string {
	return fmt.Sprintf(":%s%s%s", node.Name, node.NodeArgs, node.NodeChildren)
}

func (node *DirectiveNode) String() string {
	return node.Node()
}
