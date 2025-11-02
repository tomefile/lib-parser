package libparser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/tomefile/lib-parser/internal"
)

type Parser struct {
	Reader      *internal.AdvancedReader
	Consumer    chan *Node
	breadcrumbs []*Node
}

func New(reader *bufio.Reader, consumer chan *Node) *Parser {
	return &Parser{
		Reader:      internal.NewReader(reader),
		Consumer:    consumer,
		breadcrumbs: []*Node{},
	}
}

func (parser *Parser) ParseImmediately() {
	defer close(parser.Consumer)

	for {
		node := parser.parseNext()
		if node.IsConsumable() {
			parser.Consumer <- node
		}
		if node.IsError() {
			break
		}
	}
}

func (parser *Parser) parseNext() *Node {
	parser.Reader.RememberOffset()

	char, err := parser.Reader.Read()
	if err != nil {
		return parser.failReading(err)
	}

	switch char {

	case '}':
		if len(parser.breadcrumbs) == 0 {
			return parser.failSyntax("unexpected '}' closing a non-existant section")
		}
		parent := parser.breadcrumbs[len(parser.breadcrumbs)-1]
		parent.OffsetEnd = parser.Reader.CurrentOffset
		parser.breadcrumbs = parser.breadcrumbs[:len(parser.breadcrumbs)-1]
		return parent

	case '\n', ' ', '\t', '\v':
		return Null

	case '#':
		line, err := parser.Reader.ReadWord(internal.DelimCharset('\n'))
		if err != nil {
			return parser.failReading(err)
		}
		return parser.makeComment(line)

	case ':':
		name, err := parser.Reader.ReadWord(internal.NameCharset)
		if err != nil {
			return parser.failReading(err)
		}
		if len(name) == 0 {
			return parser.failSyntax("missing a directive name after ':'")
		}

		args, err := parser.Reader.ReadPosArgs()
		if err != nil && err != io.EOF {
			return parser.failReading(err)
		}

		char, err := parser.Reader.Read()
		if err != nil && err != io.EOF {
			return parser.failReading(err)
		}

		node := parser.makeDirective(name, args, []*Node{})

		if char == '{' {
			parser.breadcrumbs = append(parser.breadcrumbs, node)
			return Null
		} else {
			return node
		}

	default:
		name, err := parser.Reader.ReadWord(internal.NameCharset)
		if err != nil {
			return parser.failReading(err)
		}
		name = string(char) + name // We read it previously

		args, err := parser.Reader.ReadPosArgs()
		if err != nil && err != io.EOF {
			return parser.failReading(err)
		}

		if strings.HasSuffix(name, "!") {
			name = strings.TrimSuffix(name, "!")
			return parser.makeMacro(name, args)
		}

		return parser.makeExec(name, args)
	}
}

func (parser *Parser) failReading(err error) *Node {
	if err == io.EOF {
		return &Node{Type: NODE_ERROR_EOF}
	}
	return &Node{
		Type:    NODE_ERROR_READ,
		Literal: err.Error(),

		OffsetStart: parser.Reader.StoredOffset,
		OffsetEnd:   parser.Reader.CurrentOffset,
	}
}

func (parser *Parser) failSyntax(format string, a ...any) *Node {
	return &Node{
		Type:    NODE_ERROR_SYNTAX,
		Literal: fmt.Sprintf(format, a...),

		OffsetStart: parser.Reader.StoredOffset,
		OffsetEnd:   parser.Reader.CurrentOffset,
	}
}

func (parser *Parser) make(node *Node) *Node {
	if len(parser.breadcrumbs) != 0 {
		node.Parent = parser.breadcrumbs[len(parser.breadcrumbs)-1]
	}

	return node
}

func (parser *Parser) makeComment(comment string) *Node {
	return parser.make(&Node{
		Type:    NODE_COMMENT,
		Literal: comment,

		OffsetStart: parser.Reader.StoredOffset,
		OffsetEnd:   parser.Reader.CurrentOffset,
	})
}

func (parser *Parser) makeDirective(name string, args []any, children []*Node) *Node {
	return parser.make(&Node{
		Type:     NODE_DIRECTIVE,
		Literal:  name,
		Args:     args,
		Children: children,

		OffsetStart: parser.Reader.StoredOffset,
		OffsetEnd:   parser.Reader.CurrentOffset,
	})
}

func (parser *Parser) makeMacro(name string, args []any) *Node {
	return parser.make(&Node{
		Type:    NODE_MACRO,
		Literal: name,
		Args:    args,

		OffsetStart: parser.Reader.StoredOffset,
		OffsetEnd:   parser.Reader.CurrentOffset,
	})
}

func (parser *Parser) makeExec(name string, args []any) *Node {
	return parser.make(&Node{
		Type:    NODE_EXEC,
		Literal: name,
		Args:    args,

		OffsetStart: parser.Reader.StoredOffset,
		OffsetEnd:   parser.Reader.CurrentOffset,
	})
}
