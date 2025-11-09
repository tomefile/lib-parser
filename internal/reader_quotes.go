package internal

import (
	"errors"
	"strings"
)

func (reader *SourceCodeReader) ReadInsideQuotes(quote rune) (string, error) {
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
