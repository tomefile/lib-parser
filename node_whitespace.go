package libparser

type NodeWhitespace struct {
	NodeContext
}

func (node *NodeWhitespace) Context() NodeContext {
	return node.NodeContext
}

func (node *NodeWhitespace) String() string {
	return "\n"
}
