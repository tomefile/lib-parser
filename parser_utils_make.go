package libparser

func (parser *Parser) make(node_type NodeType, name string, args []any) *Node {
	return parser.fillNodeParent(parser.fillNodeOffsets(&Node{
		Type:    node_type,
		Literal: name,
		Args:    args,
	}))
}

func (parser *Parser) makeComment(comment string) *Node {
	return parser.make(NODE_COMMENT, comment, nil)
}

func (parser *Parser) makeDirective(name string, args []any) *Node {
	return parser.make(NODE_DIRECTIVE, name, args)
}

func (parser *Parser) makeMacro(name string, args []any) *Node {
	return parser.make(NODE_MACRO, name, args)
}

func (parser *Parser) makeExec(name string, args []any) *Node {
	return parser.make(NODE_EXEC, name, args)
}
