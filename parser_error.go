package libparser

import (
	"fmt"
	"io"

	liberrors "github.com/tomefile/lib-errors"
)

var EOF = &liberrors.DetailedError{Name: "EOF"}

func (parser *Parser) fillTrace(out *liberrors.DetailedError) {
	out.AddTraceItem(liberrors.TraceItem{
		Name: parser.Name,
		Col:  parser.reader.PrevCol,
		Row:  parser.reader.PrevRow,
	})
	if parser.parent != nil {
		parser.parent.fillTrace(out)
	}
}

func (parser *Parser) fail(name, details string) *liberrors.DetailedError {
	derr := &liberrors.DetailedError{
		Name:    name,
		Details: details,
		Trace:   nil,
		Context: parser.reader.ErrorContext(),
	}

	parser.fillTrace(derr)
	return derr
}

func (parser *Parser) failReading(err error) *liberrors.DetailedError {
	if err == io.EOF {
		return EOF
	}

	return parser.fail(liberrors.ERROR_READING, err.Error())
}

func (parser *Parser) failSyntax(format string, args ...any) *liberrors.DetailedError {
	return parser.fail(liberrors.ERROR_SYNTAX, fmt.Sprintf(format, args...))
}
