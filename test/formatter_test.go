package parser_test

import (
	"fmt"
	"strings"
	"testing"

	libparser "github.com/tomefile/lib-parser"
	"gotest.tools/assert"
)

const ModifierProbe = "PREFIXabc defg 123 XYZ_/hIwOrLd  (,./!@#$%^&*-=_+SUFFIX"

type TestFormat struct {
	Name           string
	ModifierResult string
	IsOptional     bool
}

func (format TestFormat) Eval(_ libparser.Scope) (string, error) {
	return format.Name, nil
}

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
				Name:       "message",
				Modifier:   nil,
				IsOptional: false,
			},
			libparser.LiteralFormat{
				Literal: " and others!",
			},
		},
	},
	{
		Input: "hi ${message:trim_suffix SUFFIX} wow!",
		Output: []libparser.FormatPart{
			libparser.LiteralFormat{
				Literal: "hi ",
			},
			libparser.VariableFormat{
				Name:       "message",
				Modifier:   libparser.GetModifier("trim_suffix", []string{"SUFFIX"}),
				IsOptional: false,
			},
			libparser.LiteralFormat{
				Literal: " wow!",
			},
		},
	},
	{
		Input: "hi ${message?:trim_suffix SUFFIX}${optional?}",
		Output: []libparser.FormatPart{
			libparser.LiteralFormat{
				Literal: "hi ",
			},
			libparser.VariableFormat{
				Name:       "message",
				Modifier:   libparser.GetModifier("trim_suffix", []string{"SUFFIX"}),
				IsOptional: true,
			},
			libparser.VariableFormat{
				Name:       "optional",
				Modifier:   nil,
				IsOptional: true,
			},
		},
	},
}

func TestFormatter(test *testing.T) {
	for i, test_case := range FormatterTestCases {
		test_case.Output = applyTestModifiers(test_case.Output)

		test.Run(
			fmt.Sprintf("%02d_with_%d_characters", i, len(test_case.Input)),
			func(test *testing.T) {
				formatter := libparser.NewFormatter(strings.NewReader(test_case.Input))

				parts, err := formatter.Parse()
				assert.NilError(test, err)

				parts = applyTestModifiers(parts)
				assert.DeepEqual(test, parts, test_case.Output)
			},
		)
	}
}

func applyTestModifiers(parts []libparser.FormatPart) []libparser.FormatPart {
	for i, part := range parts {
		switch part := part.(type) {
		case libparser.VariableFormat:
			if part.Modifier != nil {
				parts[i] = TestFormat{
					Name:           part.Name,
					ModifierResult: part.Modifier(ModifierProbe),
					IsOptional:     part.IsOptional,
				}
			}
		}
	}
	return parts
}
