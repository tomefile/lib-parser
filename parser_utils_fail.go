package libparser

import (
	"fmt"
	"io"
)

func (parser *Parser) fail(node_type NodeType, message string) *Node {
	return parser.fillNodeOffsets(&Node{
		Type:    node_type,
		Literal: message,
	})
}

func (parser *Parser) failReading(err error) *Node {
	if err == io.EOF {
		return nil
	}
	return parser.fail(NODE_ERROR_READ, err.Error())
}

func (parser *Parser) failSyntax(format string, a ...any) *Node {
	return parser.fail(NODE_ERROR_SYNTAX, fmt.Sprintf(format, a...))
}
