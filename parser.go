package libparser

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	"github.com/tomefile/lib-parser/internal"
)

type Parser struct {
	parent          *Parser
	Name            string
	reader          *internal.SourceCodeReader
	root            *NodeTree
	endOfSectionErr *DetailedError
	postProcessors  []PostProcessor
}

func New(name string, reader *bufio.Reader, processors ...PostProcessor) *Parser {
	return &Parser{
		parent: nil,
		Name:   name,
		reader: internal.NewSourceCodeReader(reader),
		root: &NodeTree{
			Tomes:        map[string]Node{},
			NodeChildren: NodeChildren{},
		},
		endOfSectionErr: nil,
		postProcessors:  processors,
	}
}

func (parser *Parser) NewNested(
	name string,
	reader *bufio.Reader,
	processors ...PostProcessor,
) *Parser {
	nested := New(name, reader, processors...)
	nested.parent = parser
	return nested
}

func (parser *Parser) Parse() (*NodeTree, *DetailedError) {
	parser.endOfSectionErr = parser.failSyntax("unexpected '}' with no matching '{' pair")

	for {
		err := parser.next(&parser.root.NodeChildren)
		if err != nil {
			if err == EOF {
				break
			}
			return parser.root, err
		}
	}
	return parser.root, nil
}

func (parser *Parser) writeNode(container *NodeChildren, node Node) (err *DetailedError) {
	for _, processor := range parser.postProcessors {
		node, err = processor(node)
		if err != nil {
			return err
		}
		if node == nil {
			// Node was discarded
			return nil
		}
	}

	*container = append(*container, node)
	return nil
}

func (parser *Parser) next(container *NodeChildren) *DetailedError {
	parser.reader.ContextReset()

	char, err := parser.reader.Read()
	if err != nil {
		return parser.failReading(err)
	}

	if unicode.IsSpace(char) {
		return nil
	}

	parser.reader.ContextBookmark()
	switch char {

	case '}':
		return parser.endOfSectionErr

	case '#':
		comment, err := parser.reader.ReadWord(internal.DelimCharset('\n'))
		if err != nil {
			return parser.failReading(err)
		}
		parser.reader.Inner.ReadRune() // Consume the \n character
		return parser.writeNode(container, &CommentNode{Contents: comment})

	case ':':
		name, err := parser.reader.ReadWord(internal.NameCharset)
		if err != nil {
			return parser.failReading(err)
		}
		if len(name) == 0 {
			return parser.failSyntax("missing a directive name after ':'")
		}

		parser.reader.ContextBookmark()
		args, err := parser.readArgs(false)
		if err != nil && err != io.EOF {
			return parser.failReading(err)
		}

		parser.reader.ContextBookmark()
		children := NodeChildren{}
		char, err := parser.reader.Peek()
		if err != nil && err != io.EOF {
			return parser.failReading(err)
		}
		if char == '{' {
			parser.reader.Read()
			children, err = parser.readChildren()
			if err != nil {
				return parser.failReading(err)
			}
		}

		return parser.writeNode(container, &DirectiveNode{
			Name:         name,
			NodeArgs:     args,
			NodeChildren: children,
		})

	default:
		name, err := parser.reader.ReadWord(internal.NameCharset)
		if err != nil {
			return parser.failReading(err)
		}
		name = string(char) + name // We read it earlier

		parser.reader.ContextBookmark()
		args, err := parser.readArgs(false)
		if err != nil && err != io.EOF {
			return parser.failReading(err)
		}

		return parser.writeNode(container, &ExecNode{
			Binary:   name,
			NodeArgs: args,
		})
	}
}

func (parser *Parser) appendArg(out NodeArgs, builder *strings.Builder) NodeArgs {
	if builder.Len() != 0 {
		out = append(out, &StringNode{Contents: builder.String()})
		builder.Reset()
	}
	return out
}

func (parser *Parser) readArgs(is_nested bool) (NodeArgs, error) {
	var builder strings.Builder
	parentheses_depth := 0
	expect_subcommand := false
	is_escaped := false
	out := NodeArgs{}

	for {
		char, err := parser.reader.Read()
		if err != nil {
			return out, err
		}

		if !is_escaped {
			if char == '\n' {
				out = parser.appendArg(out, &builder)
				return out, nil
			} else if unicode.IsSpace(char) {
				out = parser.appendArg(out, &builder)
				continue
			}
		} else {
			if char == '\n' {
				continue
			} else if unicode.IsSpace(char) {
				// NOTE: I'm not sure if this is doing anything useful
				out = parser.appendArg(out, &builder)
				continue
			}
		}

		if expect_subcommand && char != '(' {
			// It was infact not a subcommand
			expect_subcommand = false
			builder.WriteRune('$')
		}

		switch char {

		case '$':
			expect_subcommand = true

		case '(':
			parentheses_depth++
			if expect_subcommand {
				expect_subcommand = false
				out = parser.appendArg(out, &builder)
				subcommand, err := parser.readSubcommand()
				if err != nil {
					return nil, err
				}
				out = append(out, subcommand)
			}
			fallthrough

		case ')':
			parentheses_depth--
			if is_nested && parentheses_depth <= 0 {
				out = parser.appendArg(out, &builder)
				return out, io.EOF
			}

		case '{':
			parser.reader.Inner.UnreadRune()
			return parser.appendArg(out, &builder), nil

		case '\\':
			is_escaped = true

		case '"', '\'', '`':
			contents, err := parser.reader.ReadInsideQuotes(char)
			if err != nil {
				return nil, err
			}
			if char == '`' {
				out = append(out, &LiteralNode{Contents: contents})
			} else {
				out = append(out, &StringNode{Contents: contents})
			}

		case '<':
			contents, err := parser.reader.ReadInsideQuotes('>')
			if err != nil {
				return nil, err
			}
			out = append(out, &LiteralNode{Contents: "<" + contents + ">"})

		default:
			if is_escaped {
				builder.WriteRune('\\')
				is_escaped = false
			}
			builder.WriteRune(char)
		}

	}
}

func (parser *Parser) readSubcommand() (Node, error) {
	name, err := parser.reader.ReadWord(internal.NameCharset)
	if err != nil {
		return nil, err
	}
	if len(name) == 0 {
		return &StringNode{Contents: "$("}, err
	}

	parser.reader.ContextBookmark()
	args, err := parser.readArgs(true)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &ExecNode{
		Binary:   name,
		NodeArgs: args,
	}, nil
}

func (parser *Parser) readChildren() (NodeChildren, error) {
	out := NodeChildren{}

	for {
		err := parser.next(&out)
		if err != nil {
			if err == EOF || err == parser.endOfSectionErr {
				break
			}
			return out, err
		}
	}

	return out, nil
}
