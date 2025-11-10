package libparser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/tomefile/lib-parser/internal"
)

type StringFormatter struct {
	reader  *internal.SourceCodeReader
	builder strings.Builder
	out     []Segment
}

func NewStringFormatter(reader *bufio.Reader) *StringFormatter {
	return &StringFormatter{
		reader:  internal.NewSourceCodeReader(reader),
		builder: strings.Builder{},
		out:     []Segment{},
	}
}

func (formatter *StringFormatter) fail(err error) *ParsingError {
	return &ParsingError{
		Name:    ERROR_FORMATTING,
		Details: err.Error(),
		File:    "<stream>",
		Col:     formatter.reader.PrevCol,
		Row:     formatter.reader.PrevRow,
		Context: formatter.reader.GetPrintedContext(),
	}
}

func (formatter *StringFormatter) Format() ([]Segment, *ParsingError) {
	formatter.reader.ContextReset()

	for {
		char, err := formatter.reader.Read()
		if err != nil {
			if err == io.EOF {
				formatter.writeBuffer()
				return formatter.out, nil
			}
			return formatter.out, formatter.fail(err)
		}

		switch char {

		case '$':
			formatter.reader.ContextBookmark()
			formatter.writeBuffer()
			segment, err := formatter.parseVariable()
			if err != nil {
				return formatter.out, formatter.fail(err)
			}
			formatter.out = append(formatter.out, segment)

		default:
			formatter.builder.WriteRune(char)
		}
	}
}

func (formatter *StringFormatter) parseVariable() (Segment, error) {
	char, err := formatter.reader.Read()
	if err != nil {
		return nil, err
	}

	switch char {

	case '{':
		formatter.reader.ContextBookmark()
		return formatter.parseExpansion()

	default:
		formatter.reader.ContextBookmark()
		name, err := formatter.parseExpectedName(char)
		if err != nil {
			return nil, err
		}
		return &VariableSegment{
			Name:     string(char) + name,
			Modifier: nil,
		}, nil
	}
}

func (formatter *StringFormatter) parseExpansion() (Segment, error) {
	char, err := formatter.reader.Read()

	if err != nil {
		return nil, err
	}

	is_optional := false

	if char == '}' {
		// Treat empty variable expansions as a literal string
		return &LiteralNode{Contents: "${}"}, nil
	}

	formatter.reader.ContextBookmark()
	name, err := formatter.parseExpectedName(char)
	if err != nil {
		return nil, err
	}
	name = string(char) + name

	for {
		char, err = formatter.reader.Read()
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("unexpected end of file inside of variable expansion")
			}
			return nil, err
		}

		switch char {

		case '?':
			is_optional = true

		case ':':
			formatter.reader.ContextBookmark()
			format, args, err := formatter.parseFormat()
			if err != nil && err != io.EOF {
				return nil, err
			}
			return &VariableSegment{
				Name:       name,
				Modifier:   GetModifier(format, args),
				IsOptional: is_optional,
			}, nil

		case '}':
			return &VariableSegment{
				Name:       name,
				Modifier:   nil,
				IsOptional: is_optional,
			}, nil

		default:
			formatter.reader.Inner.UnreadRune()
			return nil, fmt.Errorf(
				"unexpected character %q in a variable expansion",
				char,
			)
		}
	}
}

func (formatter *StringFormatter) parseFormat() (string, []string, error) {
	formatter.reader.ContextBookmark()
	name, err := formatter.parseName()
	if err != nil {
		return "", nil, err
	}

	args := []string{}

	for {
		char, err := formatter.reader.Read()
		if err != nil {
			return name, nil, err
		}

		switch char {

		case ' ':
			formatter.reader.ContextBookmark()
			arg, err := formatter.parseName()
			if err != nil {
				return name, nil, err
			}
			args = append(args, arg)

		case '}':
			return name, args, nil

		default:
			return name, nil, fmt.Errorf(
				"unexpected character %q in a variable format",
				char,
			)
		}
	}
}

func (formatter *StringFormatter) parseExpectedName(first_char rune) (string, error) {
	if !internal.NameCharset(first_char) {
		return "", fmt.Errorf("unexpected character %q", first_char)
	}

	return formatter.parseName()
}

func (formatter *StringFormatter) parseName() (string, error) {
	var builder strings.Builder

	for {
		char, err := formatter.reader.Read()
		if err != nil {
			if err == io.EOF {
				return builder.String(), nil
			}
			return "", err
		}

		if internal.NameCharset(char) {
			builder.WriteRune(char)
		} else {
			formatter.reader.Inner.UnreadRune()
			return builder.String(), nil
		}
	}
}

func (formatter *StringFormatter) writeBuffer() {
	if formatter.builder.Len() != 0 {
		formatter.out = append(formatter.out, &LiteralNode{
			Contents: formatter.builder.String(),
		})
		formatter.builder.Reset()
	}
}
