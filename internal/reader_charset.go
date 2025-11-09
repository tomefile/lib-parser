package internal

import (
	"slices"
	"strings"
)

func (reader *SourceCodeReader) ReadWord(charset_comparator func(rune) bool) (string, error) {
	var builder strings.Builder

	for {
		char, err := reader.Read()
		if err != nil {
			return builder.String(), err
		}

		if !charset_comparator(char) {
			reader.Inner.UnreadRune()
			return builder.String(), nil
		}

		builder.WriteRune(char)
	}
}

type CharsetComparator func(rune) bool

func NameCharset(in rune) bool {
	return (in >= 'A' && in <= 'Z') ||
		(in >= 'a' && in <= 'z') ||
		in == '_' ||
		in == '-' ||
		in == '!' ||
		in == '$'
}

func DelimCharset(delims ...rune) CharsetComparator {
	return func(in rune) bool {
		return !slices.Contains(delims, in)
	}
}
