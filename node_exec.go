package libparser

import "fmt"

// An executable external program
type ExecNode struct {
	Binary string

	NodeArgs
}

func (node *ExecNode) Node() string {
	return fmt.Sprintf("%s%s", node.Binary, node.NodeArgs)
}
