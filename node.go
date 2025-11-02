package libparser

type NodeType byte

const (
	NODE_NULL NodeType = iota
	NODE_COMMENT
	NODE_DIRECTIVE
	NODE_MACRO
	NODE_EXEC

	NODE_ERROR_READ
	NODE_ERROR_SYNTAX
	NODE_ERROR_EOF

	ErrMin = NODE_ERROR_READ
	ErrMax = NODE_ERROR_EOF
)

var Null = &Node{Type: NODE_NULL}

type Node struct {
	Parent *Node // Optional

	Type    NodeType
	Literal string
	Args    []any // can either contain [string] or [*Node]

	// The offset in bytes where this node begins
	OffsetStart uint
	// The offset in bytes where this node ends
	OffsetEnd uint

	Children []*Node // Optional
}

func (node *Node) IsConsumable() bool {
	return node.Type != NODE_NULL && node.Type != NODE_ERROR_EOF
}

func (node *Node) IsError() bool {
	return node.Type >= ErrMin && node.Type <= ErrMax
}
