package libparser

import (
	"strings"

	liberrors "github.com/tomefile/lib-errors"
)

type Hook func(original Node) (modified Node, derr *liberrors.DetailedError)

func ExcludeHook[T Node](node Node) (Node, *liberrors.DetailedError) {
	switch node.(type) {
	case T:
		return nil, nil
	}
	return node, nil
}

func NoShebangHook(node Node) (Node, *liberrors.DetailedError) {
	switch node := node.(type) {
	case *NodeComment:
		if strings.HasPrefix(node.Contents, "!") {
			return nil, nil
		}
	}
	return node, nil
}

// FIXME: Make sure it works.
func StreamHook(channel chan Node) Hook {
	return func(node Node) (Node, *liberrors.DetailedError) {
		channel <- node
		return node, nil
	}
}
