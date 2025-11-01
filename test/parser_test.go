package parser_test

import (
	"bufio"
	"bytes"
	"fmt"
	"sync"
	"testing"

	libparser "github.com/tomefile/lib-parser"
)

const ExampleProgram1 = `# Example program, привет мир 🟡!
:include <std>

echo "Hello World!" \
	"and another line" \
	"and another."
`

func TestParser(test *testing.T) {
	var wg sync.WaitGroup
	var buffer bytes.Buffer

	buffer.WriteString(ExampleProgram1)

	consumer := make(chan libparser.Statement)
	parser := libparser.New(bufio.NewReader(&buffer), consumer)
	wg.Go(parser.Parse)

	for statement := range consumer {
		fmt.Printf(
			"==> %#v\n -> %q\n\n",
			statement,
			string([]byte(ExampleProgram1)[statement.OffsetStart:statement.OffsetEnd]),
		)

		switch statement.Kind {
		case libparser.SK_NULL:
		case libparser.SK_COMMENT:
		case libparser.SK_DIRECTIVE:
		case libparser.SK_EXEC:
		case libparser.SK_MACRO:
		case libparser.SK_READ_ERROR:
		case libparser.SK_SYNTAX_ERROR:
		default:
			test.Errorf("unexpected libparser.StatementKind: %#v", statement.Kind)
		}
	}
}
