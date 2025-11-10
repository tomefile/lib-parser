package libparser

import "fmt"

// A literal string does not get modified in any way.
type LiteralNode struct {
	Contents string
}

func (node *LiteralNode) Node() string {
	return fmt.Sprintf("[literal %q]", node.Contents)
}

func (node *LiteralNode) Eval(_ Locals) (string, error) {
	return node.Contents, nil
}
