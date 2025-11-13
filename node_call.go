package libparser

import "fmt"

// An executable external program
type CallNode struct {
	Macro string

	NodeArgs
}

func (node *CallNode) Node() string {
	return fmt.Sprintf("%s!%s", node.Macro, node.NodeArgs)
}

func (node *CallNode) String() string {
	return node.Node()
}
