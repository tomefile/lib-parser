package readers

import (
	"errors"
	"slices"
	"strings"
)

type CharsetComparator func(rune) bool

func NameCharset(in rune) bool {
	return (in >= 'A' && in <= 'Z') ||
		(in >= 'a' && in <= 'z') ||
		(in >= '0' && in <= '9') ||
		in == '_' ||
		in == '-' ||
		in == '!' ||
		in == '$'
}

func FilenameCharset(in rune) bool {
	return NameCharset(in) ||
		in == '.' ||
		in == '~' ||
		in == '/' ||
		in == '\\'
}

func QuotesCharset(in rune) bool {
	return in == '\'' ||
		in == '"' ||
		in == '`'
}

func (reader *Reader) ReadSequence(comparator CharsetComparator) (string, error) {
	var builder strings.Builder

	for {
		char, err := reader.Read()
		if err != nil {
			return builder.String(), err
		}

		if !comparator(char) {
			reader.Unread()
			return builder.String(), nil
		}

		builder.WriteRune(char)
	}
}

func (reader *Reader) ReadDelimited(delims ...rune) (string, error) {
	var builder strings.Builder

	for {
		char, err := reader.Read()
		if err != nil {
			return builder.String(), err
		}

		if slices.Contains(delims, char) {
			reader.Unread()
			return builder.String(), nil
		}

		builder.WriteRune(char)
	}
}

func (reader *Reader) ReadInsideQuotes(quote rune) (string, error) {
	var builder strings.Builder
	is_escaped := false

	for {
		char, err := reader.Read()
		if err != nil {
			return builder.String(), err
		}

		switch char {

		case quote:
			if !is_escaped {
				return builder.String(), nil
			}

		case '\\':
			is_escaped = true

		case '\n':
			if !is_escaped {
				return builder.String(), errors.New("unexpected new line inside of quotes")
			}
		}

		is_escaped = false
		builder.WriteRune(char)
	}
}
