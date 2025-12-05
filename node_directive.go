package libparser

type NodeDirective struct {
	Name string
	NodeContext
	NodeArgs
	NodeChildren
}

func (node *NodeDirective) Context() NodeContext {
	return node.NodeContext
}

func (node *NodeDirective) String() string {
	return ":" +
		node.Name +
		node.NodeArgs.String() +
		node.NodeChildren.String()
}
