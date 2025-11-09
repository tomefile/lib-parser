package libparser

type Node interface {
	Node() string
}

type NodeArgs []Node

type NodeChildren []Node

type NodeTree struct {
	Tomes map[string]Node

	NodeChildren
}
