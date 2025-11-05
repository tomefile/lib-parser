package libparser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/tomefile/lib-parser/internal"
)

type Formatter struct {
	Reader  *bufio.Reader
	builder strings.Builder
	out     []FormatPart
}

func NewFormatter(reader io.Reader) *Formatter {
	return &Formatter{
		Reader:  bufio.NewReader(reader),
		builder: strings.Builder{},
		out:     []FormatPart{},
	}
}

func (formatter *Formatter) Read() (rune, error) {
	char, _, err := formatter.Reader.ReadRune()
	return char, err
}

func (formatter *Formatter) Parse() ([]FormatPart, error) {
	for {
		char, err := formatter.Read()
		if err != nil {
			if err == io.EOF {
				formatter.popLiteral()
				return formatter.out, nil
			}
			return formatter.out, err
		}

		switch char {

		case '$':
			formatter.popLiteral()
			part, err := formatter.parseVariable()
			if err != nil {
				return formatter.out, err
			}
			formatter.out = append(formatter.out, part)

		default:
			formatter.builder.WriteRune(char)

		}
	}
}

func (formatter *Formatter) parseVariable() (FormatPart, error) {
	char, err := formatter.Read()
	if err != nil {
		return nil, err
	}

	switch char {

	case '{':
		return formatter.parseExpansion()

	default:
		name, err := formatter.parseExpectedName(char)
		if err != nil {
			return nil, err
		}
		return VariableFormat{
			Name:     string(char) + name,
			Modifier: nil,
		}, nil
	}
}

func (formatter *Formatter) parseExpansion() (FormatPart, error) {
	char, err := formatter.Read()
	if err != nil {
		return nil, err
	}

	is_optional := false

	if char == '}' {
		// Treat empty variable expansions as a literal string
		return LiteralFormat{Literal: "${}"}, nil
	}

	name, err := formatter.parseExpectedName(char)
	if err != nil {
		return nil, err
	}
	name = string(char) + name

	for {
		char, err = formatter.Read()
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
			format, args, err := formatter.parseFormat()
			if err != nil && err != io.EOF {
				return nil, err
			}
			return VariableFormat{
				Name:       name,
				Modifier:   GetModifier(format, args),
				IsOptional: is_optional,
			}, nil

		case '}':
			return VariableFormat{
				Name:       name,
				Modifier:   nil,
				IsOptional: is_optional,
			}, nil

		default:
			formatter.Reader.UnreadRune()
			return nil, fmt.Errorf(
				"unexpected character %q in a variable expansion",
				char,
			)
		}
	}
}

func (formatter *Formatter) parseFormat() (string, []string, error) {
	name, err := formatter.parseName()
	if err != nil {
		return "", nil, err
	}

	args := []string{}

	for {
		char, err := formatter.Read()
		if err != nil {
			return name, nil, err
		}

		switch char {

		case ' ':
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

func (formatter *Formatter) parseExpectedName(first_char rune) (string, error) {
	if !internal.NameCharset(first_char) {
		return "", fmt.Errorf("unexpected character %q", first_char)
	}

	return formatter.parseName()
}

func (formatter *Formatter) parseName() (string, error) {
	var builder strings.Builder

	for {
		char, err := formatter.Read()
		if err != nil {
			if err == io.EOF {
				return builder.String(), nil
			}
			return "", err
		}

		if internal.NameCharset(char) {
			builder.WriteRune(char)
		} else {
			formatter.Reader.UnreadRune()
			return builder.String(), nil
		}
	}
}

func (formatter *Formatter) popLiteral() {
	if formatter.builder.Len() != 0 {
		formatter.out = append(formatter.out, LiteralFormat{
			Literal: formatter.builder.String(),
		})
		formatter.builder.Reset()
	}
}
