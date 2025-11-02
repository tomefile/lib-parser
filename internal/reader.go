package internal

import (
	"bufio"
	"errors"
	"io"
	"slices"
	"strings"
)

// Buffered + tracked reader
type AdvancedReader struct {
	Inner *bufio.Reader

	StoredOffset  uint
	CurrentOffset uint

	StoredCol, StoredRow uint
	PrevCol, PrevRow     uint
	Col, Row             uint
}

func NewReader(reader *bufio.Reader) *AdvancedReader {
	return &AdvancedReader{
		Inner:         reader,
		StoredOffset:  1,
		CurrentOffset: 1,
		Col:           1,
		Row:           1,
	}
}

func (reader *AdvancedReader) RememberOffset() *AdvancedReader {
	reader.StoredOffset = reader.CurrentOffset
	return reader
}

func (reader *AdvancedReader) RememberPosition() *AdvancedReader {
	reader.StoredCol = reader.Col
	reader.StoredRow = reader.Row
	return reader
}

func (reader *AdvancedReader) Peek() (byte, error) {
	data, err := reader.Inner.Peek(1)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

// Reads the next UTF-8 rune
func (reader *AdvancedReader) Read() (rune, error) {
	reader.PrevRow = reader.Row
	reader.PrevCol = reader.Col

	char, size, err := reader.Inner.ReadRune()
	if err != nil {
		return 0, err
	}

	reader.CurrentOffset += uint(size)

	if char == '\n' {
		reader.Row++
		reader.Col = 0
	} else {
		reader.Col += uint(size)
	}

	return char, nil
}

func (reader *AdvancedReader) ReadWord(charset_matcher CharsetMatcher) (string, error) {
	var builder strings.Builder

	for {
		char, err := reader.Read()
		if err != nil {
			return builder.String(), err
		}

		if !charset_matcher(char) {
			return builder.String(), nil
		}

		builder.WriteRune(char)
	}
}

type CharsetMatcher func(rune) bool

func NameCharset(in rune) bool {
	return (in >= 'A' && in <= 'Z') ||
		(in >= 'a' && in <= 'z') ||
		in == '_' ||
		in == '-' ||
		in == '!' ||
		in == '$'
}

func DelimCharset(delims ...rune) CharsetMatcher {
	return func(in rune) bool {
		return !slices.Contains(delims, in)
	}
}

func (reader *AdvancedReader) ReadInsideQuotes(quote rune) (string, error) {
	var builder strings.Builder

	for {
		char, err := reader.Read()
		if err != nil {
			return builder.String(), err
		}

		switch char {
		case quote:
			return builder.String(), nil
		case '\n':
			return builder.String(), errors.New("unexpected new line inside of quotes")
		}

		builder.WriteRune(char)
	}
}

func (reader *AdvancedReader) ReadPosArgs(is_nested bool) ([]any, error) {
	var builder strings.Builder
	out := []any{}
	is_escaped := false
	expect_subcommand := false
	parentheses_depth := 0

	for {
		// region: Temporary fix
		peek_char, err := reader.Peek()
		if err != nil {
			return out, err
		}
		if peek_char == '{' {
			if builder.Len() != 0 {
				out = append(out, builder.String())
			}
			return out, nil
		}
		// endregion

		char, err := reader.Read()
		if err != nil {
			return out, err
		}

		if expect_subcommand && char != '(' {
			expect_subcommand = false
			builder.WriteRune('$')
		}

		switch char {

		case '\n':
			if !is_escaped {
				if builder.Len() != 0 {
					out = append(out, builder.String())
				}
				return out, nil
			}
			is_escaped = false

		case '\\':
			is_escaped = true

		case ' ', '\t':
			if builder.Len() != 0 {
				out = append(out, builder.String())
				builder.Reset()
			}

		case '"', '`', '\'':
			contents, err := reader.ReadInsideQuotes(char)
			if err != nil {
				return nil, err
			}
			out = append(out, contents)

		case '$':
			expect_subcommand = true

		case '(':
			parentheses_depth++
			if expect_subcommand {
				subcommand, err := reader.ReadSubcommand()
				if err != nil {
					return nil, err
				}
				out = append(out, subcommand)
				expect_subcommand = false
				continue
			}
			fallthrough

		case ')':
			parentheses_depth--
			if is_nested && parentheses_depth <= 0 {
				if builder.Len() != 0 {
					out = append(out, builder.String())
				}
				return out, io.EOF
			}

		default:
			if is_escaped {
				builder.WriteRune('\\')
			}
			is_escaped = false

			builder.WriteRune(char)
		}
	}
}

type Subcommand struct {
	Name    string
	Args    []any
	IsMacro bool
}

func (reader *AdvancedReader) ReadSubcommand() (*Subcommand, error) {
	subcommand := &Subcommand{}

	name, err := reader.ReadWord(NameCharset)
	if err != nil {
		return nil, err
	}
	subcommand.Name = strings.TrimSuffix(name, "!")
	subcommand.IsMacro = strings.HasSuffix(name, "!")

	args, err := reader.ReadPosArgs(true)
	if err != nil && err != io.EOF {
		return nil, err
	}
	subcommand.Args = args

	return subcommand, nil
}
