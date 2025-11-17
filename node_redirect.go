package libparser

import "fmt"

type RedirectType byte

const (
	REDIRECT_STDOUT RedirectType = iota
	REDIRECT_STDERR
	REDIRECT_STDIN
	REDIRECT_HEREDOC
	REDIRECT_HERESTR
)

type RedirectNode struct {
	Type   RedirectType
	Source Node
	Dest   Node
}

func (node *RedirectNode) Node() string {
	switch node.Type {
	case REDIRECT_STDOUT:
		// [source] %s [dest]
		return fmt.Sprintf(
			"%s %s %s",
			node.Source.Node(),
			node.Type.String(),
			node.Dest.Node(),
		)
	default:
		// [dest] %s [source]
		return fmt.Sprintf(
			"%s %s %s",
			node.Dest.Node(),
			node.Type.String(),
			node.Source.Node(),
		)
	}
}

func (node *RedirectNode) String() string {
	return node.Node()
}

func (node RedirectType) String() string {
	switch node {
	case REDIRECT_HEREDOC:
		return "<<"
	case REDIRECT_HERESTR:
		return "<<<"
	case REDIRECT_STDIN:
		return "<"
	case REDIRECT_STDOUT:
		return ">"
	}

	return ""
}

type ChildRedirectNode struct {
	Source  Node
	OutDest Node
	ErrDest Node
}

func (node *ChildRedirectNode) Node() string {
	return fmt.Sprintf(
		"%s > %s >> %s",
		node.Source.Node(),
		node.OutDest.Node(),
		node.ErrDest.Node(),
	)
}

func (node *ChildRedirectNode) String() string {
	return node.Node()
}
