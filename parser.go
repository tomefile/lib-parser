package libparser

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	liberrors "github.com/tomefile/lib-errors"
	"github.com/tomefile/lib-parser/internal"
)

type Parser struct {
	parent          *Parser
	Name            string
	reader          *internal.SourceCodeReader
	root            *NodeTree
	endOfSectionErr *liberrors.DetailedError
	PostProcessors  []PostProcessor
}

func New(file File) *Parser {
	return &Parser{
		parent: nil,
		Name:   file.Name(),
		reader: internal.NewSourceCodeReader(bufio.NewReader(file)),
		root: &NodeTree{
			Tomes:        map[string]Node{},
			NodeChildren: NodeChildren{},
		},
		endOfSectionErr: nil,
		PostProcessors:  []PostProcessor{},
	}
}

// Appends the [PostProcessor] to be applied to every single node before it gets appended to the tree.
//
// NOTE: Order matters (sequentially from first to last)
func (parser *Parser) With(processor PostProcessor) *Parser {
	parser.PostProcessors = append(parser.PostProcessors, processor)
	return parser
}

// Used for error tracing
func (parser *Parser) SetParent(parent *Parser) *Parser {
	parser.parent = parent
	return parser
}

func (parser *Parser) Parse() (*NodeTree, *liberrors.DetailedError) {
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

func (parser *Parser) writeNode(
	container *NodeChildren,
	node Node,
) (derr *liberrors.DetailedError) {
	// The reason it's calculated early is because it can be changed
	// during post-processing, but it is still a tome.
	tome_name := ""
	switch node := node.(type) {
	case *DirectiveNode:
		if node.Name == "tome" && len(node.NodeArgs) > 0 {
			arg := node.NodeArgs[0]
			tome_name = arg.Node()
		}
	}

	for _, processor := range parser.PostProcessors {
		node, derr = processor(node)
		if derr != nil {
			derr.Context = parser.reader.ErrorContext()
			// Highlight the entire buffer
			derr.Context.Highlighted = strings.TrimSuffix(derr.Context.Buffer, "\n")
			derr.Context.Buffer = ""
			parser.fillTrace(derr)
			return derr
		}
		if node == nil {
			// Node was discarded
			return nil
		}
	}

	if tome_name != "" {
		parser.root.Tomes[tome_name] = node
	}

	*container = append(*container, node)
	return nil
}

func (parser *Parser) next(container *NodeChildren) *liberrors.DetailedError {
	parser.reader.ContextReset()

	char, err := parser.reader.Read()
	if err != nil {
		return parser.failReading(err)
	}

	if unicode.IsSpace(char) || char == ';' {
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
		if !internal.NameCharset(char) {
			return parser.failSyntax(
				"unexpected %q at the beginning of the line. If it's a part of a statement on the line above, append '\\' at the end of the previous line.",
				char,
			)
		}

		parser.reader.Inner.UnreadRune()
		node, derr := parser.readExec()
		if derr != nil {
			return derr
		}
		return parser.writeNode(container, node)
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
			if char == '\n' || char == ';' {
				out = parser.appendArg(out, &builder)
				return out, nil
			} else if unicode.IsSpace(char) {
				out = parser.appendArg(out, &builder)
				continue
			}
		} else {
			if char == '\n' || char == ';' {
				continue
			} else if unicode.IsSpace(char) {
				// NOTE: I'm not sure if this is doing anything useful
				out = parser.appendArg(out, &builder)
				continue
			}
		}

		if expect_subcommand && char != '(' && char != '{' {
			// It was infact not a subcommand
			expect_subcommand = false
			builder.WriteRune('$')
		}

		switch char {

		case '|':
			parser.reader.Inner.UnreadRune()
			out = parser.appendArg(out, &builder)
			return out, nil

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
			if !expect_subcommand {
				parser.reader.Inner.UnreadRune()
				return parser.appendArg(out, &builder), nil
			}
			contents, err := parser.reader.ReadInsideQuotes('}')
			if err != nil {
				return nil, err
			}
			builder.WriteString("${" + contents + "}")
			expect_subcommand = false

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

func (parser *Parser) skipWhitespace() *liberrors.DetailedError {
	for {
		char, err := parser.reader.Read()
		if err != nil {
			return parser.failReading(err)
		}

		if !unicode.IsSpace(char) {
			parser.reader.Inner.UnreadRune()
			return nil
		}
	}
}

func (parser *Parser) readExec() (Node, *liberrors.DetailedError) {
	if derr := parser.skipWhitespace(); derr != nil {
		return nil, derr
	}

	name, err := parser.reader.ReadWord(internal.NameCharset)
	if err != nil {
		return nil, parser.failReading(err)
	}

	parser.reader.ContextBookmark()
	args, err := parser.readArgs(false)
	if err != nil && err != io.EOF {
		return nil, parser.failReading(err)
	}

	var node Node
	if strings.HasSuffix(name, "!") {
		node = &CallNode{
			Macro:    name[:len(name)-1],
			NodeArgs: args,
		}
	} else {
		node = &ExecNode{
			Binary:   name,
			NodeArgs: args,
		}
	}

	peek_char, _ := parser.reader.Peek()
	if peek_char == '|' {
		parser.reader.Read()

		dest_node, dest_err := parser.readExec()
		if dest_err != nil {
			return nil, dest_err
		}

		return &PipeNode{
			Source: node,
			Dest:   dest_node,
		}, nil
	}

	return node, nil
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
