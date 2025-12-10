package libparser

import (
	"strings"

	liberrors "github.com/tomefile/lib-errors"
)

// A function to be called on every Node before the parser attempts to append it to the tree.
//
// Return `nil` to discard the Node.
type Hook func(original Node) (modified Node, derr *liberrors.DetailedError)

// Discards Nodes of type T
func ExcludeHook[T Node](node Node) (Node, *liberrors.DetailedError) {
	switch node.(type) {
	case T:
		return nil, nil
	}
	return node, nil
}

// Removes the UNIX shebang
func NoShebangHook(node Node) (Node, *liberrors.DetailedError) {
	switch node := node.(type) {
	case *NodeComment:
		if strings.HasPrefix(node.Contents, "!") {
			return nil, nil
		}
	}
	return node, nil
}

// Puts the node to chan before it gets appended to the tree.
// Useful for tracking the progress of parsing in very large files.
//
// # Example usage:
//
//	channel := make(chan libparser.Node)
//
//	parser := libparser.New(file)
//	parser.Hooks = []libparser.Hook{
//	    libparser.StreamHook(channel),
//	}
//
//	go func() {
//	    defer close(channel) // IMPORTANT! Without this you'll hang the process FOREVER.
//	    if derr := parser.Run(); derr != nil {
//	        ...
//	    }
//	}()
//
//	for node := range channel {
//	    ...
//	}
func StreamHook(channel chan Node) Hook {
	return func(node Node) (Node, *liberrors.DetailedError) {
		channel <- node
		return node, nil
	}
}
