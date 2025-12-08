package libparser

import (
	"fmt"
	"io"

	liberrors "github.com/tomefile/lib-errors"
)

var EOF = &liberrors.DetailedError{Name: "EOF", Details: "End of File"}

var UNEXPECTED_EOF = &liberrors.DetailedError{Name: "EOF", Details: "Unexpected End of File"}

var EOB = &liberrors.DetailedError{Name: "EOB", Details: "End of Block"}

var EOA = &liberrors.DetailedError{Name: "EOA", Details: "End of Arguments"}

func (parser *Parser) fillErrorTrace(derr *liberrors.DetailedError) {
	derr.AddTraceItem(liberrors.TraceItem{
		Name: parser.File.Name(),
		Col:  parser.reader.PrevCol,
		Row:  parser.reader.PrevRow,
	})
	if parser.Parent != nil {
		parser.Parent.fillErrorTrace(derr)
	}
}

func (parser *Parser) fail(at uint, name, details string) *liberrors.DetailedError {
	derr := &liberrors.DetailedError{
		Name:    name,
		Details: details,
		Trace:   nil,
		Context: parser.reader.Context(at),
	}

	parser.fillErrorTrace(derr)
	return derr
}

func (parser *Parser) failReading(err error) *liberrors.DetailedError {
	if err == io.EOF {
		return EOF
	}

	if derr, ok := err.(*liberrors.DetailedError); ok {
		return derr
	}

	return parser.fail(parser.reader.Offset-1, liberrors.ERROR_READING, err.Error())
}

func (parser *Parser) failSyntax(at uint, format string, args ...any) *liberrors.DetailedError {
	return parser.fail(at, liberrors.ERROR_SYNTAX, fmt.Sprintf(format, args...))
}

func (parser *Parser) failSyntaxHere(format string, args ...any) *liberrors.DetailedError {
	return parser.failSyntax(parser.reader.Offset-1, format, args...)
}
