package libparser

func (parser *Parser) fillNodeOffsets(node *Node) *Node {
	node.OffsetStart = parser.Reader.StoredOffset
	node.OffsetEnd = parser.Reader.CurrentOffset
	node.Col = parser.Reader.PrevCol
	node.Row = parser.Reader.PrevRow
	return node
}

func (parser *Parser) fillNodeParent(node *Node) *Node {
	if len(parser.breadcrumbs) != 0 {
		node.Parent = parser.parentNode()
	}
	return node
}
