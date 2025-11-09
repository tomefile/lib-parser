package libparser

import "fmt"

// A string that can have variable expansions inside.
// Uses backticks (`) instead of quotes.
type StringNode struct {
	Contents string
}

func (node *StringNode) Node() string {
	return fmt.Sprintf("[string %q]", node.Contents)
}
