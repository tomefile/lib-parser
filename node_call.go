package libparser

import "fmt"

// An executable external program
type CallNode struct {
	Macro string

	NodeArgs
}

func (node *CallNode) Node() string {
	return fmt.Sprintf("[call %q]", node.Macro)
}
