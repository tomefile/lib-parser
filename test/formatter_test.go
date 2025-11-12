package libparser_test

import (
	"fmt"
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

func (format TestFormat) Eval(_ libparser.Locals) (string, error) {
	return format.Name, nil
}

type FormatterTestCase struct {
	Input  string
	Output []libparser.Segment
}

var FormatterTestCases = []FormatterTestCase{
	{
		Input: "hello $message and others!",
		Output: []libparser.Segment{
			&libparser.LiteralNode{
				Contents: "hello ",
			},
			&libparser.VariableSegment{
				Name:       "message",
				Modifier:   nil,
				IsOptional: false,
			},
			&libparser.LiteralNode{
				Contents: " and others!",
			},
		},
	},
	{
		Input: "hi ${message:trim_suffix SUFFIX} wow!",
		Output: []libparser.Segment{
			&libparser.LiteralNode{
				Contents: "hi ",
			},
			&libparser.VariableSegment{
				Name:       "message",
				Modifier:   libparser.GetModifier("trim_suffix", []string{"SUFFIX"}),
				IsOptional: false,
			},
			&libparser.LiteralNode{
				Contents: " wow!",
			},
		},
	},
	{
		Input: "hi ${message?:trim_suffix SUFFIX}${optional?}",
		Output: []libparser.Segment{
			&libparser.LiteralNode{
				Contents: "hi ",
			},
			&libparser.VariableSegment{
				Name:       "message",
				Modifier:   libparser.GetModifier("trim_suffix", []string{"SUFFIX"}),
				IsOptional: true,
			},
			&libparser.VariableSegment{
				Name:       "optional",
				Modifier:   nil,
				IsOptional: true,
			},
		},
	},
	{
		Input: "stay $ unmodified ${}",
		Output: []libparser.Segment{
			&libparser.LiteralNode{
				Contents: "stay $ unmodified ",
			},
			&libparser.LiteralNode{
				Contents: "${}",
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
				formatter := libparser.NewStringFormatter(test_case.Input)

				parts, derr := formatter.Format()
				if derr != nil {
					derr.Print(test.Output())
					test.FailNow()
				}

				parts = applyTestModifiers(parts)
				assert.DeepEqual(test, parts, test_case.Output)
			},
		)
	}
}

func applyTestModifiers(parts []libparser.Segment) []libparser.Segment {
	for i, part := range parts {
		switch part := part.(type) {
		case *libparser.VariableSegment:
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
