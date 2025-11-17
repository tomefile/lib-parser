package libparser

import (
	"fmt"
	"io"

	liberrors "github.com/tomefile/lib-errors"
)

// End of File
var EOF = &liberrors.DetailedError{Name: "EOF"}

// Unexpected End of File
var UNEXPECTED_EOF = &liberrors.DetailedError{
	Name:    liberrors.ERROR_READING,
	Details: "unexpected EOF when expecting a '}'",
}

// End of Section
var EOS = &liberrors.DetailedError{Name: "EOS"}

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
		Context: parser.reader.Context(),
	}

	parser.fillTrace(derr)
	return derr
}

func (parser *Parser) failReading(err error) *liberrors.DetailedError {
	if err == io.EOF {
		return EOF
	}

	if derr, ok := err.(*liberrors.DetailedError); ok {
		return derr
	}

	return parser.fail(liberrors.ERROR_READING, err.Error())
}

func (parser *Parser) failSyntax(format string, args ...any) *liberrors.DetailedError {
	return parser.fail(liberrors.ERROR_SYNTAX, fmt.Sprintf(format, args...))
}
