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

	ErrMin = NODE_ERROR_READ
	ErrMax = NODE_ERROR_SYNTAX
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
	Col, Row  uint

	Children []*Node // Optional
}

func (node *Node) IsError() bool {
	return node.Type >= ErrMin && node.Type <= ErrMax
}
