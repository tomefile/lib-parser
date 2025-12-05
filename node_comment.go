package libparser

import "fmt"

type NodeComment struct {
	Contents string
	NodeContext
}

func (node *NodeComment) Context() NodeContext {
	return node.NodeContext
}

func (node *NodeComment) String() string {
	return fmt.Sprintf("#%s", node.Contents)
}
