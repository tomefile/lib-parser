package libparser

import (
	"strings"

	liberrors "github.com/tomefile/lib-errors"
)

type PostProcessor func(Node) (Node, *liberrors.DetailedError)

// Discards nodes of type [T] from the tree
func PostExclude[T Node](node Node) (Node, *liberrors.DetailedError) {
	switch node.(type) {
	case T:
		return nil, nil
	}
	return node, nil
}

// Discards shebang (unix) comment
func PostNoShebang(node Node) (Node, *liberrors.DetailedError) {
	switch node := node.(type) {
	case *CommentNode:
		if strings.HasPrefix(node.Contents, "!") {
			return nil, nil
		}
	}
	return node, nil
}
