package libparser

import (
	"fmt"
	"io"
	"strings"

	libescapes "github.com/bbfh-dev/lib-ansi-escapes"
)

const (
	ERROR_READING    = "Reading Error"
	ERROR_SYNTAX     = "Syntax Error"
	ERROR_INTERNAL   = "Internal Error"
	ERROR_FORMATTING = "Formatting Error"
)

var EOF = &DetailedError{Name: "EOF"}

type DetailedError struct {
	Name    string
	Details string

	File     string
	Col, Row uint
	Context  string
}

func (err *DetailedError) Error() string {
	return fmt.Sprintf("(%s) %s", err.Name, err.Details)
}

func (err *DetailedError) BeautyPrint(writer io.Writer) {
	// TODO: Add trace
	fmt.Fprintf(
		writer,
		"[!] %s\n    in %s:%d:%d\n\n%s\n\n[?] Details\n    %s%s%s",
		err.Name,
		err.File,
		err.Row,
		err.Col,
		err.Context,
		libescapes.TextColorBrightRed,
		err.Details,
		libescapes.ColorReset,
	)
}

func (err *DetailedError) GetBeautyPrinted() string {
	var builder strings.Builder
	err.BeautyPrint(&builder)
	return builder.String()
}

func (parser *Parser) fail(name, details string) *DetailedError {
	return &DetailedError{
		Name:    name,
		Details: details,
		File:    parser.Name,
		Col:     parser.reader.PrevCol,
		Row:     parser.reader.PrevRow,
		Context: parser.reader.GetPrintedContext(),
	}
}

func (parser *Parser) failReading(err error) *DetailedError {
	if err == io.EOF {
		return EOF
	}

	return parser.fail(ERROR_READING, err.Error())
}

func (parser *Parser) failSyntax(format string, args ...any) *DetailedError {
	return parser.fail(ERROR_SYNTAX, fmt.Sprintf(format, args...))
}
