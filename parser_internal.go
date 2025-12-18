package libparser

import (
	liberrors "github.com/tomefile/lib-errors"
)

func (parser *Parser) makeContext(offset uint) NodeContext {
	return NodeContext{
		OffsetStart: offset,
		OffsetEnd:   parser.reader.Offset,
	}
}

func (parser *Parser) process(node Node) (Node, *liberrors.DetailedError) {
	// The reason it's calculated early is because it can be changed
	// during post-processing, but it is still a tome.
	// NOTE: This has a side-effect of tomes not being discarded by hooks.
	tome_name := ""
	switch node := node.(type) {
	case *NodeDirective:
		if node.Name == "tome" && len(node.NodeArgs) > 0 {
			arg := node.NodeArgs[0]
			switch arg := arg.(type) {
			case *NodeLiteral:
				tome_name = arg.Contents
			case *NodeString:
				tome_name = arg.Segments.String()
			default:
				tome_name = arg.String()
			}
		}
	}

	// Run hooks
	offset_start := node.Context().OffsetStart
	var derr *liberrors.DetailedError
	for _, hook := range parser.Hooks {
		node, derr = hook(node)
		if derr != nil {
			derr.Context = parser.reader.Context(offset_start)
			parser.fillErrorTrace(derr)
			return node, derr
		}
		if node == nil {
			// Node was discarded
			return nil, nil
		}
	}

	if tome_name != "" {
		parser.Result.Tomes[tome_name] = node.(*NodeDirective)
	}

	return node, nil
}

func (parser *Parser) write(node Node) (derr *liberrors.DetailedError) {
	node, derr = parser.process(node)
	if derr != nil || node == nil {
		return derr
	}

	*parser.container = append(*parser.container, node)
	return nil
}

func (parser *Parser) escaped(char, comp rune) bool {
	return parser.reader.Previous() == '\\' && char == comp
}
