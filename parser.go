package libparser

import (
	"bufio"
	"io"
	"strings"

	liberrors "github.com/tomefile/lib-errors"
	"github.com/tomefile/lib-parser/readers"
)

type Locals map[string]string

type Parser struct {
	Parent    *Parser
	File      File
	Result    *NodeRoot
	Hooks     []Hook
	reader    *readers.Reader
	container *NodeChildren
}

func New(file File) *Parser {
	root := &NodeRoot{
		Tomes: map[string]Node{},
		NodeContext: NodeContext{
			OffsetStart: 0,
			OffsetEnd:   0,
		},
		NodeChildren: NodeChildren{},
	}
	return &Parser{
		Parent:    nil,
		File:      file,
		Result:    root,
		Hooks:     []Hook{},
		reader:    readers.New(bufio.NewReader(file)),
		container: &root.NodeChildren,
	}
}

func (parser *Parser) Run() *liberrors.DetailedError {
	for {
		derr := parser.next()
		switch derr {

		case nil:
			continue

		case EOF:
			return nil

		case UNEXPECTED_EOF:
			return parser.failReading(derr)

		case EOB:
			return parser.failSyntaxHere("unexpected '}' with no matching '{' pair")

		default:
			return derr
		}
	}
}

func (parser *Parser) next() *liberrors.DetailedError {
	start_offset := parser.reader.Offset

	char, err := parser.reader.Read()
	if err != nil {
		return parser.failReading(err)
	}

	if readers.WhitespaceCharset(char) || char == ';' {
		if char == '\n' {
			return parser.write(&NodeWhitespace{NodeContext: parser.makeContext(start_offset)})
		}
		return nil
	}

	switch char {

	case '}':
		char, _ := parser.reader.Read()
		if char != '\n' {
			parser.reader.Unread()
		}
		return EOB

	case '#':
		comment, err := parser.reader.ReadDelimited(true, '\n')
		if err != nil {
			return parser.failReading(err)
		}
		return parser.write(&NodeComment{
			Contents:    comment,
			NodeContext: parser.makeContext(start_offset),
		})

	case ':':
		name, err := parser.reader.ReadSequence(readers.NameCharset)
		if err != nil {
			return parser.failReading(err)
		}
		if len(name) == 0 {
			return parser.failSyntaxHere("missing a directive name after ':'")
		}

		args, derr := parser.readArgs()
		if derr != nil && derr != EOF {
			return derr
		}

		children, err := parser.readChildren()
		if err != nil && err != io.EOF {
			return parser.failReading(err)
		}

		return parser.write(&NodeDirective{
			Name:         name,
			NodeArgs:     args,
			NodeChildren: children,
			NodeContext:  parser.makeContext(start_offset),
		})
	}

	if !readers.FilenameCharset(char) && !readers.QuotesCharset(char) {
		return parser.failSyntaxHere(
			"unexpected %q at the start of a statement. If it belongs to the statement above, add a '\\' to the end of the previous line.",
			char,
		)
	}

	parser.reader.Unread()
	node, derr := parser.readStatement()
	if derr != nil {
		return derr
	}

	peek, _ := parser.reader.Peek()
	if peek == '>' || peek == '<' {
		redirection, derr := parser.readRedirection()
		if derr != nil {
			return derr
		}
		redirection.Source = node
		redirection.NodeContext = parser.makeContext(start_offset)
		return parser.write(redirection)
	}

	return parser.write(node)
}

func (parser *Parser) readRedirection() (*NodeRedirect, *liberrors.DetailedError) {
	out := &NodeRedirect{}

	for {
		char, err := parser.reader.Read()
		if err != nil {
			return out, parser.failReading(err)
		}

		switch char {
		case '\n':
			return out, nil

		case '<':
			filename, derr := parser.readFilename()
			if derr != nil {
				return out, derr
			}
			out.Stdin = filename

		case '>':
			char, err := parser.reader.Peek()
			if err != nil {
				return nil, parser.failReading(err)
			}
			switch char {
			case '>':
				parser.reader.Read()
				filename, derr := parser.readFilename()
				if derr != nil {
					return out, derr
				}
				out.Stderr = filename
			default:
				filename, derr := parser.readFilename()
				if derr != nil {
					return out, derr
				}
				out.Stdout = filename
			}
		}
	}
}

func (parser *Parser) readFilename() (*NodeString, *liberrors.DetailedError) {
	// TODO: Add string parsing
	char, err := parser.reader.Read()
	if err != nil {
		return nil, parser.failReading(err)
	}

	switch char {

	case ' ':
		return parser.readFilename()

	case '\'', '"', '`':
		contents, err := parser.reader.ReadInsideQuotes(char)
		if err != nil {
			return nil, parser.failReading(err)
		}
		return NewSimpleNodeString(contents), nil

	default:
		parser.reader.Unread()
		filename, err := parser.reader.ReadSequence(readers.FilenameCharset)
		if err != nil {
			return nil, parser.failReading(err)
		}
		return NewSimpleNodeString(filename), nil
	}
}

func (parser *Parser) readStatement() (Node, *liberrors.DetailedError) {
	start_offset := parser.reader.Offset

	string_name, err := parser.readFilename()
	if err != nil {
		return nil, parser.failReading(err)
	}
	name := string_name.String()

	args, derr := parser.readArgs()
	if derr != nil && derr != EOF {
		return nil, derr
	}

	var node Node
	if strings.HasSuffix(name, "!") {
		node = &NodeCall{
			Macro:       name[:len(name)-1],
			NodeArgs:    args,
			NodeContext: parser.makeContext(start_offset),
		}
	} else {
		node = &NodeExec{
			Name:        name,
			NodeArgs:    args,
			NodeContext: parser.makeContext(start_offset),
		}
	}

	char, _ := parser.reader.Peek()

	switch char {

	case '|':
		parser.reader.Read()
		parser.reader.ReadSequence(readers.WhitespaceCharset)

		target, derr := parser.readStatement()
		if derr != nil {
			return nil, derr
		}
		return &NodePipe{
			Source:      node,
			Dest:        target,
			NodeContext: parser.makeContext(start_offset),
		}, nil

	case ')':
		parser.reader.Read()
	}

	return node, nil
}

func (parser *Parser) readArgs() (NodeArgs, *liberrors.DetailedError) {
	var out = NodeArgs{}
	for {
		arg, derr := parser.readArg()
		if arg != nil {
			out = append(out, arg)
		}
		if derr != nil {
			if derr == EOA {
				return out, nil
			}
			return out, derr
		}
	}
}

func (parser *Parser) readArg() (Node, *liberrors.DetailedError) {
	var start_offset = parser.reader.Offset
	var out = SegmentedString{}
	var current_segment strings.Builder

	for {
		char, err := parser.reader.Read()
		if err != nil {
			return nil, parser.failReading(err)
		}

		if readers.ArglistTeminatingCharset(char) {
			if current_segment.Len() != 0 {
				out = append(out, &LiteralStringSegment{Contents: current_segment.String()})
			}

			if parser.escaped(char, '\n') {
				return &NodeWhitespace{IsLineBreak: true}, nil
			}

			if char == ' ' || char == '\t' || parser.escaped(char, '\n') {
				if len(out) == 0 {
					return nil, nil
				}
				return &NodeString{
					Segments:    out,
					NodeContext: parser.makeContext(start_offset),
				}, nil
			}

			if current_segment.Len() == 0 && len(out) == 0 {
				if char != '\n' {
					parser.reader.Unread()
				}
				return nil, EOA
			}

			var derr *liberrors.DetailedError
			if char == '\n' || char == ';' || char == ')' {
				derr = EOA
			}

			return &NodeString{
				Segments:    out,
				NodeContext: parser.makeContext(start_offset),
			}, derr
		}

		if char != '$' {
			if char == '\'' || char == '"' || char == '`' {
				contents, err := parser.reader.ReadInsideQuotes(char)
				if err != nil {
					return nil, parser.failReading(err)
				}
				literal := &NodeLiteral{
					Contents:    contents,
					NodeContext: parser.makeContext(start_offset),
				}
				if char == '\'' {
					return literal, nil
				}
				return literal.ToStringNode(), nil
			}
			// was the previous character '\\'
			if parser.escaped(0, 0) {
				current_segment.WriteRune('\\')
			}
			if char != '\\' {
				current_segment.WriteRune(char)
			}
			continue
		}

		if current_segment.Len() != 0 {
			out = append(out, &LiteralStringSegment{Contents: current_segment.String()})
			current_segment.Reset()
		}

		char_after_dollar, err := parser.reader.Read()
		if err != nil {
			return nil, parser.failReading(err)
		}

		if readers.WhitespaceCharset(char_after_dollar) || char_after_dollar == ';' {
			literal := &NodeLiteral{
				Contents:    "$",
				NodeContext: parser.makeContext(start_offset),
			}
			return literal.ToStringNode(), nil
		}

		switch char_after_dollar {

		case '(':
			node, derr := parser.readStatement()
			if derr != nil {
				return nil, derr
			}
			return node, nil

		case '{':
			segment, derr := parser.readVariableExpansion()
			if derr != nil {
				return nil, derr
			}
			out = append(out, segment)

		default:
			word, err := parser.reader.ReadSequence(readers.NameCharset)
			if err != nil {
				return nil, parser.failReading(err)
			}
			out = append(out, &VariableStringSegment{
				Name:       string(char_after_dollar) + word,
				Modifiers:  []StringModifier{},
				IsOptional: false,
			})
		}
	}
}

func (parser *Parser) readChildren() (NodeChildren, error) {
	out := NodeChildren{}
	backup := parser.container
	parser.container = &out
	defer func() {
		parser.container = backup
	}()

	char, err := parser.reader.Read()
	if err != nil {
		return out, err
	}
	if char != '{' {
		parser.reader.Unread()
		return out, nil
	}

	char, err = parser.reader.Read()
	if err != nil {
		return out, err
	}

	for {
		derr := parser.next()
		switch derr {

		case nil:
			continue

		case EOB:
			return out, nil

		case EOF:
			return out, UNEXPECTED_EOF

		default:
			return out, derr
		}
	}
}

func (parser *Parser) readVariableExpansion() (*VariableStringSegment, *liberrors.DetailedError) {
	name, err := parser.reader.ReadSequence(readers.NameCharset)
	if err != nil {
		return nil, parser.failReading(err)
	}

	char, err := parser.reader.Read()
	if err != nil {
		return nil, parser.failReading(err)
	}

	var optional bool
	if char == '?' {
		optional = true

		char, err = parser.reader.Read()
		if err != nil {
			return nil, parser.failReading(err)
		}
	}

	if char == '}' {
		return &VariableStringSegment{
			Name:       name,
			Modifiers:  []StringModifier{},
			IsOptional: optional,
		}, nil
	}

	if char != ':' {
		parser.reader.Unread()
		return nil, parser.failSyntaxHere("unexpected %q in a variable expansion", char)
	}

	modifiers := []StringModifier{}

	for {
		modifier, derr := parser.readVariableModifier()
		if modifier.Name != "" {
			modifiers = append(modifiers, modifier)
		}
		if derr != nil {
			if derr == EOA {
				return &VariableStringSegment{
					Name:       name,
					Modifiers:  SortNotModifierToEnd(modifiers),
					IsOptional: optional,
				}, nil
			}
			return nil, derr
		}
	}
}

func (parser *Parser) readVariableModifier() (StringModifier, *liberrors.DetailedError) {
	offset_start := parser.reader.Offset

	modifier_name, err := parser.reader.ReadSequence(readers.NameCharset)
	if err != nil {
		return StringModifier{}, parser.failReading(err)
	}

	char, err := parser.reader.Read()
	if err != nil {
		return StringModifier{}, parser.failReading(err)
	}

	args := []*NodeString{}

	switch char {

	case ' ':
	read_argument:
		arg, derr := parser.readFilename()
		if derr != nil {
			return StringModifier{}, derr
		}
		if len(arg.Segments) != 0 {
			args = append(args, arg)
		}

		next_char, err := parser.reader.Read()
		if err != nil {
			return StringModifier{}, parser.failReading(err)
		}
		if next_char == '}' || next_char == ':' {
			modifier, err := GetModifier(ModifierName(modifier_name), args)
			if err != nil {
				return StringModifier{}, parser.fail(
					offset_start,
					liberrors.ERROR_VALIDATION,
					err.Error(),
				)
			}
			if next_char == '}' {
				return modifier, EOA
			}
			return modifier, nil
		}
		if next_char == ' ' {
			goto read_argument
		}
		goto unexpected_character

	case '}', ':':
		modifier, err := GetModifier(ModifierName(modifier_name), []*NodeString{})
		if err != nil {
			return StringModifier{}, parser.fail(
				offset_start,
				liberrors.ERROR_VALIDATION,
				err.Error(),
			)
		}
		if char == '}' {
			return modifier, EOA
		}
		return modifier, nil
	}

unexpected_character:
	return StringModifier{}, parser.failSyntaxHere(
		"unexpected character %q in a variable expansion modifier",
		char,
	)
}
