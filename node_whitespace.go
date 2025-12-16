package libparser

type NodeWhitespace struct {
	IsLineBreak bool
	NodeContext
}

func (node *NodeWhitespace) Context() NodeContext {
	return node.NodeContext
}

func (node *NodeWhitespace) String() string {
	if node.IsLineBreak {
		return "\\\n"
	}
	return "\n"
}
