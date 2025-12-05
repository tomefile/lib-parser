package libparser

import "fmt"

type NodeLiteral struct {
	Contents string
	NodeContext
}

func (node *NodeLiteral) Context() NodeContext {
	return node.NodeContext
}

func (node *NodeLiteral) String() string {
	return fmt.Sprintf("'%s'", node.Contents)
}

func (node *NodeLiteral) ToStringNode() *NodeString {
	return &NodeString{
		Segments: SegmentedString{
			&LiteralStringSegment{
				Contents: node.Contents,
			},
		},
		NodeContext: node.Context(),
	}
}
