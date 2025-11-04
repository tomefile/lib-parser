package libparser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/tomefile/lib-parser/internal"
)

type ArgPart struct {
	Literal    string
	Format     func(string) string
	IsVariable bool
}

func ParseArg(arg any) ([]ArgPart, error) {
	str_arg, ok := arg.(string)
	if !ok {
		return nil, fmt.Errorf("ParseArg() expected a string but got %q instead", arg)
	}

	var builder strings.Builder
	out := []ArgPart{}
	reader := bufio.NewReader(strings.NewReader(str_arg))

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				if builder.Len() != 0 {
					out = append(out, newLiteralPart(&builder))
					builder.Reset()
				}
				return out, nil
			}
			return out, err
		}

		switch char {

		case '$':
			if builder.Len() != 0 {
				out = append(out, newLiteralPart(&builder))
				builder.Reset()
			}
			part, err := readVariable(reader)
			if err != nil {
				return out, err
			}
			out = append(out, part)

		default:
			builder.WriteRune(char)

		}
	}
}

func readVariable(reader *bufio.Reader) (ArgPart, error) {
	char, _, err := reader.ReadRune()
	if err != nil {
		return ArgPart{}, err
	}

	switch char {

	case '{':
		return readVariableExpansion(reader)

	default:
		if !internal.NameCharset(char) {
			return ArgPart{}, fmt.Errorf("unexpected character %q in name", char)
		}
		var builder strings.Builder
		builder.WriteRune(char)

		for {
			char, _, err := reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					return ArgPart{
						Literal:    builder.String(),
						Format:     nil,
						IsVariable: true,
					}, nil
				}
				return ArgPart{}, err
			}

			if internal.NameCharset(char) {
				builder.WriteRune(char)
			} else {
				reader.UnreadRune()
				return ArgPart{
					Literal:    builder.String(),
					Format:     nil,
					IsVariable: true,
				}, nil
			}
		}
	}
}

func readVariableExpansion(reader *bufio.Reader) (ArgPart, error) {
	return ArgPart{}, nil
}

func newLiteralPart(builder *strings.Builder) ArgPart {
	return ArgPart{
		Literal:    builder.String(),
		Format:     nil,
		IsVariable: false,
	}
}
