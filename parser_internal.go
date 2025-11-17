package libparser

import (
	"io"
	"strings"

	liberrors "github.com/tomefile/lib-errors"
	"github.com/tomefile/lib-parser/readers"
)

func (parser *Parser) process(node Node) (Node, *liberrors.DetailedError) {
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

	var derr *liberrors.DetailedError
	for _, processor := range parser.PostProcessors {
		node, derr = processor(node)
		if derr != nil {
			derr.Context = parser.reader.Context()
			// Highlight the entire buffer
			derr.Context.Highlighted = strings.TrimSuffix(derr.Context.Buffer, "\n")
			derr.Context.Buffer = ""
			parser.fillTrace(derr)
			return node, derr
		}
		if node == nil {
			// Node was discarded
			return nil, nil
		}
	}

	if tome_name != "" {
		parser.root.Tomes[tome_name] = node
	}

	return node, nil
}

func (parser *Parser) write(container *NodeChildren, node Node) (derr *liberrors.DetailedError) {
	node, derr = parser.process(node)
	if derr != nil || node == nil {
		return derr
	}

	*container = append(*container, node)
	return nil
}

func (parser *Parser) next(container *NodeChildren) *liberrors.DetailedError {
	parser.reader.MarkContext()

	char, err := parser.reader.Read()
	if err != nil {
		return parser.failReading(err)
	}

	if readers.WhitespaceCharset(char) || char == ';' {
		return nil
	}

	parser.reader.MarkSegment()
	switch char {

	case '}':
		return EOS

	case '#':
		comment, err := parser.reader.ReadDelimited('\n')
		if err != nil {
			return parser.failReading(err)
		}
		parser.reader.Inner.ReadRune() // Consume the \n character
		return parser.write(container, &CommentNode{Contents: comment})

	case ':':
		name, err := parser.reader.ReadSequence(readers.NameCharset)
		if err != nil {
			return parser.failReading(err)
		}
		if len(name) == 0 {
			return parser.failSyntax("missing a directive name after ':'")
		}

		parser.reader.MarkSegment()
		args, err := parser.readArgs(false)
		if err != nil && err != io.EOF {
			return parser.failReading(err)
		}

		parser.reader.MarkSegment()
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

		return parser.write(container, &DirectiveNode{
			Name:         name,
			NodeArgs:     args,
			NodeChildren: children,
		})

	default:
		if !readers.NameCharset(char) {
			return parser.failSyntax(
				"unexpected %q at the beginning of the line. If it's a part of a statement on the line above, append '\\' at the end of the previous line.",
				char,
			)
		}

		parser.reader.Unread()
		node, derr := parser.readExec()
		if derr != nil {
			return derr
		}
		return parser.write(container, node)
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
			} else if readers.WhitespaceCharset(char) {
				out = parser.appendArg(out, &builder)
				continue
			}
		} else {
			if char == '\n' || char == ';' {
				continue
			} else if readers.WhitespaceCharset(char) {
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

		case '|', '>', '<':
			parser.reader.Unread()
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
				parser.reader.Unread()
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
	name, err := parser.reader.ReadSequence(readers.NameCharset)
	if err != nil {
		return nil, err
	}
	if len(name) == 0 {
		return &StringNode{Contents: "$("}, err
	}

	parser.reader.MarkSegment()
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

		if !readers.WhitespaceCharset(char) {
			parser.reader.Unread()
			return nil
		}
	}
}

func (parser *Parser) readRedirection(node Node) (Node, *liberrors.DetailedError) {
	if derr := parser.skipWhitespace(); derr != nil {
		return nil, derr
	}

	char, err := parser.reader.Read()
	if err != nil {
		return nil, parser.failReading(err)
	}

	switch char {
	case '>':
		peek_char, _ := parser.reader.Peek()
		redirect_type := REDIRECT_STDOUT
		if peek_char == '>' {
			parser.reader.Read()
			redirect_type = REDIRECT_STDERR
		}

		dest_node, derr := parser.readFilename()
		if derr != nil {
			return nil, derr
		}

		following, derr := parser.readRedirection(node)
		if derr != nil {
			return nil, derr
		}
		if following == nil {
			return &RedirectNode{
				Type:   redirect_type,
				Source: node,
				Dest:   dest_node,
			}, nil
		}
		return &ChildRedirectNode{
			Source:  node,
			OutDest: dest_node,
			ErrDest: following,
		}, nil

	case '<':
		peek_char, _ := parser.reader.Peek()
		redirect_type := REDIRECT_STDIN
		if peek_char == '<' {
			parser.reader.Read()
			redirect_type = REDIRECT_HEREDOC
			peek_char, _ = parser.reader.Peek()
			if peek_char == '<' {
				parser.reader.Read()
				redirect_type = REDIRECT_HERESTR
			}
		}

		src_node, derr := parser.readFilename()
		if derr != nil {
			return nil, derr
		}

		return &RedirectNode{
			Type:   redirect_type,
			Source: src_node,
			Dest:   node,
		}, nil

	default:
		return nil, nil
	}
}

func (parser *Parser) readFilename() (Node, *liberrors.DetailedError) {
	if derr := parser.skipWhitespace(); derr != nil {
		return nil, derr
	}

	var builder strings.Builder

	char, err := parser.reader.Read()
	if err != nil {
		return nil, parser.failReading(err)
	}

	switch {

	case readers.FilenameCharset(char):
		builder.WriteRune(char)
		word, err := parser.reader.ReadDelimited('\n', ' ')
		if err != nil {
			return nil, parser.failReading(err)
		}
		builder.WriteString(word)
		return &StringNode{Contents: builder.String()}, nil

	case readers.QuotesCharset(char):
		contents, err := parser.reader.ReadInsideQuotes(char)
		if err != nil {
			return nil, parser.failReading(err)
		}
		if char == '`' {
			return &LiteralNode{Contents: contents}, nil
		}
		return &StringNode{Contents: contents}, nil

	default:
		parser.reader.Unread()
		return nil, parser.failSyntax("unexpected character %q in file name", char)
	}
}

func (parser *Parser) readExec() (Node, *liberrors.DetailedError) {
	if derr := parser.skipWhitespace(); derr != nil {
		return nil, derr
	}

	name, err := parser.reader.ReadSequence(readers.NameCharset)
	if err != nil {
		return nil, parser.failReading(err)
	}

	parser.reader.MarkSegment()
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
	switch peek_char {

	case '|':
		parser.reader.Read()

		dest_node, dest_err := parser.readExec()
		if dest_err != nil {
			return nil, dest_err
		}

		return &PipeNode{
			Source: node,
			Dest:   dest_node,
		}, nil

	case '>', '<':
		re_node, derr := parser.readRedirection(node)
		if derr != nil {
			return nil, derr
		}
		if re_node == nil {
			return node, nil
		}
		return re_node, nil

	}

	return node, nil
}

func (parser *Parser) readChildren() (NodeChildren, error) {
	out := NodeChildren{}

	for {
		derr := parser.next(&out)
		if derr != nil {
			if derr == EOS {
				break
			}
			if derr == EOF {
				return out, UNEXPECTED_EOF
			}
			return out, derr
		}
	}

	return out, nil
}
