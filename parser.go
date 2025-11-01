package libparser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/tomefile/lib-parser/internal"
)

type Parser struct {
	Reader   *internal.AdvancedReader
	Consumer chan Statement
	parents  []*Statement
}

func New(reader *bufio.Reader, consumer chan Statement) *Parser {
	return &Parser{
		Reader:   internal.NewReader(reader),
		Consumer: consumer,
		parents:  []*Statement{},
	}
}

func (parser *Parser) Parse() {
	defer close(parser.Consumer)

	for {
		statement := parser.parseStatement()
		if statement.IsConsumable() {
			parser.Consumer <- statement
		}
		if statement.IsError() {
			break
		}
	}
}

func (parser *Parser) parseStatement() Statement {
	parser.Reader.RememberOffset()

	char, err := parser.Reader.Read()
	if err != nil {
		return parser.failReading(err)
	}

	switch char {

	case '}':
		if len(parser.parents) == 0 {
			return parser.failSyntax("unexpected '}' closing a non-existant section")
		}
		parent := parser.parents[len(parser.parents)-1]
		parent.OffsetEnd = parser.Reader.CurrentOffset
		parser.parents = parser.parents[:len(parser.parents)-1]
		return *parent

	case '\n', ' ', '\t', '\v':
		return NullStatement

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

		statement := parser.makeDirective(name, args, []Statement{})

		if char == '{' {
			parser.parents = append(parser.parents, &statement)
			return NullStatement
		} else {
			return statement
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

func (parser *Parser) failReading(err error) Statement {
	if err == io.EOF {
		return Statement{Kind: SK_EOF_ERROR}
	}
	return Statement{
		Kind:    SK_READ_ERROR,
		Literal: err.Error(),

		OffsetStart: parser.Reader.StoredOffset,
		OffsetEnd:   parser.Reader.CurrentOffset,
	}
}

func (parser *Parser) failSyntax(format string, a ...any) Statement {
	return Statement{
		Kind:    SK_SYNTAX_ERROR,
		Literal: fmt.Sprintf(format, a...),

		OffsetStart: parser.Reader.StoredOffset,
		OffsetEnd:   parser.Reader.CurrentOffset,
	}
}

func (parser *Parser) make(statement Statement) Statement {
	if len(parser.parents) != 0 {
		parser.parents[len(parser.parents)-1].Children = append(
			parser.parents[len(parser.parents)-1].Children,
			statement,
		)
		return NullStatement
	}

	return statement
}

func (parser *Parser) makeComment(comment string) Statement {
	return parser.make(Statement{
		Kind:    SK_COMMENT,
		Literal: comment,

		OffsetStart: parser.Reader.StoredOffset,
		OffsetEnd:   parser.Reader.CurrentOffset,
	})
}

func (parser *Parser) makeDirective(name string, args []string, children []Statement) Statement {
	return parser.make(Statement{
		Kind:     SK_DIRECTIVE,
		Literal:  name,
		Args:     args,
		Children: children,

		OffsetStart: parser.Reader.StoredOffset,
		OffsetEnd:   parser.Reader.CurrentOffset,
	})
}

func (parser *Parser) makeMacro(name string, args []string) Statement {
	return parser.make(Statement{
		Kind:    SK_MACRO,
		Literal: name,
		Args:    args,

		OffsetStart: parser.Reader.StoredOffset,
		OffsetEnd:   parser.Reader.CurrentOffset,
	})
}

func (parser *Parser) makeExec(name string, args []string) Statement {
	return parser.make(Statement{
		Kind:    SK_EXEC,
		Literal: name,
		Args:    args,

		OffsetStart: parser.Reader.StoredOffset,
		OffsetEnd:   parser.Reader.CurrentOffset,
	})
}
