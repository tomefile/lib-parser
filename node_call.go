package libparser

type NodeCall struct {
	Macro string
	NodeContext
	NodeArgs
}

func (node *NodeCall) Context() NodeContext {
	return node.NodeContext
}

func (node *NodeCall) String() string {
	return node.Macro + "!" +
		node.NodeArgs.String()
}
