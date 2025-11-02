package libparser

import (
	"bufio"
	"io"
	"strings"

	"github.com/tomefile/lib-parser/internal"
)

type Parser struct {
	Reader  *internal.AdvancedReader
	Channel chan *Node

	breadcrumbs []*Node
}

func New(reader *bufio.Reader, consumer chan *Node) *Parser {
	return &Parser{
		Reader:      internal.NewReader(reader),
		Channel:     consumer,
		breadcrumbs: []*Node{},
	}
}

// Nodes are returned as soon as they are read (fast).
// Sets [Node.Parent].
// [Node.Children] is always nil
func (parser *Parser) ParseFragmented() {
	defer close(parser.Channel)

	for {
		node := parser.next()
		if node == nil { // EOF
			return
		}

		if node == Null {
			continue
		}

		parser.Channel <- node
		if node.IsError() {
			return
		}
	}
}

// Nodes are returned only once all of their children are read (slow).
// Sets [Node.Children] if applicable.
// [Node.Parent] is always nil to prevent pointer recursion
func (parser *Parser) ParseComplete() {
	defer close(parser.Channel)

	var buffered *Node

	for {
		node := parser.next()
		if node == nil { // EOF
			if buffered != nil {
				parser.Channel <- buffered
			}
			return
		}

		if node == Null {
			continue
		}

		if node.IsError() {
			if buffered != nil {
				parser.Channel <- buffered
			}
			parser.Channel <- node
			return
		}

		if node.Parent == nil {
			if buffered != nil {
				parser.Channel <- buffered
			}
			buffered = node
		} else {
			node.Parent.Children = append(node.Parent.Children, node)
			node.Parent = nil
		}
	}
}

func (parser *Parser) next() *Node {
	parser.Reader.RememberOffset()

	char, err := parser.Reader.Read()
	if err != nil {
		return parser.failReading(err)
	}

	switch char {

	case '\n', ' ', '\t', '\v':
		return Null

	case '#':
		comment, err := parser.Reader.ReadWord(internal.DelimCharset('\n'))
		if err != nil {
			return parser.failReading(err)
		}
		return parser.makeComment(comment)

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

		node := parser.makeDirective(name, args)

		if char == '{' {
			parser.breadcrumbs = append(parser.breadcrumbs, node)
		}
		return node

	case '}':
		if len(parser.breadcrumbs) == 0 {
			return parser.failSyntax("unexpected '}' closing a non-existant section")
		}
		parser.parentNode().OffsetEnd = parser.Reader.CurrentOffset
		parser.breadcrumbs = parser.breadcrumbs[:len(parser.breadcrumbs)-1]
		return Null

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

func (parser *Parser) parentNode() *Node {
	if len(parser.breadcrumbs) == 0 {
		return nil
	}
	return parser.breadcrumbs[len(parser.breadcrumbs)-1]
}
