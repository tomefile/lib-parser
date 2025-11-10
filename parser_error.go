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

type TraceItem struct {
	Name     string
	Col, Row uint
}

type DetailedError struct {
	Name    string
	Details string

	Trace   []TraceItem
	Context string
}

func (err *DetailedError) Error() string {
	return fmt.Sprintf("(%s) %s", err.Name, err.Details)
}

func (err *DetailedError) BeautyPrint(writer io.Writer) {
	fmt.Fprintf(writer, "[!] %s\n", err.Name)

	if len(err.Trace) > 0 {
		fmt.Fprintf(
			writer,
			"    in %s:%d:%d\n",
			err.Trace[0].Name,
			err.Trace[0].Row,
			err.Trace[0].Col,
		)

		for _, item := range err.Trace[1:] {
			fmt.Fprintf(
				writer,
				"    └─ from %s:%d:%d\n",
				item.Name,
				item.Row,
				item.Col,
			)
		}
	}

	fmt.Fprintf(
		writer,
		"\n%s\n\n[?] Details\n    %s%s%s",
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
	trace := []TraceItem{
		{
			Name: parser.Name,
			Col:  parser.reader.PrevCol,
			Row:  parser.reader.PrevRow,
		},
	}

	return &DetailedError{
		Name:    name,
		Details: details,
		Trace:   trace,
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
