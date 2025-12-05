package libparser

type NodeExec struct {
	Name string
	NodeContext
	NodeArgs
}

func (node *NodeExec) Context() NodeContext {
	return node.NodeContext
}

func (node *NodeExec) String() string {
	return node.Name +
		node.NodeArgs.String()
}
