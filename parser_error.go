package libparser

import (
	"fmt"
	"io"

	libescapes "github.com/bbfh-dev/lib-ansi-escapes"
)

const (
	ERROR_READING    = "Reading Error"
	ERROR_SYNTAX     = "Syntax Error"
	ERROR_INTERNAL   = "Internal Error"
	ERROR_FORMATTING = "Formatting Error"
)

var EOF = &ParsingError{Name: "EOF"}

type ParsingError struct {
	Name    string
	Details string

	File     string
	Col, Row uint
	Context  string
}

func (err *ParsingError) Error() string {
	return fmt.Sprintf("(%s) %s", err.Name, err.Details)
}

func (err *ParsingError) BeautyPrint(writer io.Writer) {
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

func (parser *Parser) fail(name, details string) *ParsingError {
	return &ParsingError{
		Name:    name,
		Details: details,
		File:    parser.Name,
		Col:     parser.reader.PrevCol,
		Row:     parser.reader.PrevRow,
		Context: parser.reader.GetPrintedContext(),
	}
}

func (parser *Parser) failReading(err error) *ParsingError {
	if err == io.EOF {
		return EOF
	}

	return parser.fail(ERROR_READING, err.Error())
}

func (parser *Parser) failSyntax(format string, args ...any) *ParsingError {
	return parser.fail(ERROR_SYNTAX, fmt.Sprintf(format, args...))
}
