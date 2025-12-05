package libparser

type NodePipe struct {
	Source Node
	Dest   Node
	NodeContext
}

func (node *NodePipe) Context() NodeContext {
	return node.NodeContext
}

func (node *NodePipe) String() string {
	return node.Source.String() + " | " +
		node.Dest.String()
}
