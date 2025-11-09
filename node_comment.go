package libparser

import "fmt"

// Can be used for documentation
type CommentNode struct {
	Contents string
}

func (node *CommentNode) Node() string {
	return fmt.Sprintf("[comment %q]", node.Contents)
}
