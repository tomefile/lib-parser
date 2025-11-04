package parser_test

import (
	"strings"
	"testing"

	libparser "github.com/tomefile/lib-parser"
	"gotest.tools/assert"
)

func TestFormatter(test *testing.T) {
	formatter := libparser.NewFormatter(strings.NewReader("hello $message and others!"))
	parts, err := formatter.Parse()
	assert.NilError(test, err)
	assert.DeepEqual(test, parts, []libparser.FormatPart{
		libparser.LiteralFormat{
			Literal: "hello ",
		},
		libparser.VariableFormat{
			Name:     "message",
			Modifier: nil,
		},
		libparser.LiteralFormat{
			Literal: " and others!",
		},
	})
}
