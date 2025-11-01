package internal

import (
	"bufio"
	"slices"
	"strings"
)

// Buffered + tracked reader
type AdvancedReader struct {
	Inner *bufio.Reader

	StoredOffset  uint
	CurrentOffset uint

	Col uint
	Row uint
}

func NewReader(reader *bufio.Reader) *AdvancedReader {
	return &AdvancedReader{
		Inner:         reader,
		StoredOffset:  0,
		CurrentOffset: 0,
		Col:           0,
		Row:           0,
	}
}

func (reader *AdvancedReader) RememberOffset() *AdvancedReader {
	reader.StoredOffset = reader.CurrentOffset
	return reader
}

// Reads the next UTF-8 rune
func (reader *AdvancedReader) Read() (rune, error) {
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

func (reader *AdvancedReader) ReadPosArgs() ([]string, error) {
	var builder strings.Builder
	out := []string{}
	is_escaped := false

	for {
		char, err := reader.Read()
		if err != nil {
			return out, err
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
			// FIXME: Quotes don't have to be closed to be valid. This isn't good
			contents, err := reader.ReadWord(DelimCharset(char, '\n'))
			if err != nil {
				return nil, err
			}
			out = append(out, contents)

		default:
			if is_escaped {
				builder.WriteRune('\\')
			}
			is_escaped = false

			builder.WriteRune(char)
		}
	}
}
