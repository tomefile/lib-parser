package parser_test

import (
	"fmt"
	"strings"
	"testing"

	libparser "github.com/tomefile/lib-parser"
	"gotest.tools/assert"
)

type FormatterTestCase struct {
	Input  string
	Output []libparser.FormatPart
}

var FormatterTestCases = []FormatterTestCase{
	{
		Input: "hello $message and others!",
		Output: []libparser.FormatPart{
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
		},
	},
}

func TestFormatter(test *testing.T) {
	for i, test_case := range FormatterTestCases {
		test.Run(
			fmt.Sprintf("%02d_with_%d_characters", i, len(test_case.Input)),
			func(test *testing.T) {
				formatter := libparser.NewFormatter(strings.NewReader(test_case.Input))

				parts, err := formatter.Parse()
				assert.NilError(test, err)

				assert.DeepEqual(test, parts, test_case.Output)
			},
		)
	}
}
