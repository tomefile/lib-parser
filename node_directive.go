package libparser

import "fmt"

type DirectiveNode struct {
	Name string

	NodeArgs
	NodeChildren
}

func (node *DirectiveNode) Node() string {
	return fmt.Sprintf("[directive %q]", node.Name)
}
