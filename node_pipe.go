package libparser

import "fmt"

// Can either contain *ExecNode or *CallNode
type PipeNode struct {
	Source Node
	Dest   Node
}

func (node *PipeNode) Node() string {
	return fmt.Sprintf("%s | %s", node.Source.Node(), node.Dest.Node())
}

func (node *PipeNode) String() string {
	return node.Node()
}
