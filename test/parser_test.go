package parser_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	libparser "github.com/tomefile/lib-parser"
	"gotest.tools/assert"
)

var NoChildren = []libparser.Statement{}

var FileTestCases = map[string][]libparser.Statement{
	"01_basic.tome": {
		{
			Kind:    libparser.SK_COMMENT,
			Literal: " Example program, привет мир 🟡!",
		},
		{
			Kind:     libparser.SK_DIRECTIVE,
			Literal:  "include",
			Args:     []string{"<std>"},
			Children: NoChildren,
		},
		{
			Kind:    libparser.SK_EXEC,
			Literal: "echo",
			Args:    []string{"Hello World!", "and another line", "and another."},
		},
	},
	"02_directive_body.tome": {
		{
			Kind:    libparser.SK_EXEC,
			Literal: "echo",
			Args:    []string{"1"},
		},
		{
			Kind:    libparser.SK_DIRECTIVE,
			Literal: "section",
			Args:    []string{"Hello World!"},
			Children: []libparser.Statement{
				{
					Kind:    libparser.SK_EXEC,
					Literal: "echo",
					Args:    []string{"1.1"},
				},
				{
					Kind:    libparser.SK_EXEC,
					Literal: "echo",
					Args:    []string{"1.2"},
				},
			},
		},
	},
}

func TestAll(test *testing.T) {
	dir := "data"

	entries, err := os.ReadDir(dir)
	assert.NilError(test, err)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tome") {
			continue
		}

		file_path := filepath.Join(dir, entry.Name())

		contents, err := os.ReadFile(file_path)
		assert.NilError(test, err)

		file, err := os.OpenFile(file_path, os.O_RDONLY, os.ModePerm)
		assert.NilError(test, err)
		defer file.Close()

		testFile(test, file, entry.Name(), contents)
	}
}

func testFile(test *testing.T, file *os.File, name string, buffer []byte) {
	test.Run(name, func(test *testing.T) {
		test_case, exists := FileTestCases[name]
		if !exists {
			test.Fatalf("missing test case for %q", name)
		}

		consumer := make(chan libparser.Statement)
		parser := libparser.New(bufio.NewReader(file), consumer)

		go parser.Parse()

		var i int
		for statement := range consumer {
			if len(test_case) <= i {
				test.Fatalf("unknown statement: %#v", statement)
			}

			data, err := json.MarshalIndent(statement, "", strings.Repeat(" ", 4))
			assert.NilError(test, err)
			test.Logf(
				"%s\nRaw data: %q",
				string(data),
				string(buffer[statement.OffsetStart:statement.OffsetEnd]),
			)

			// Using individual fields because not every field is important
			assert.DeepEqual(test, test_case[i].Kind, statement.Kind)
			assert.DeepEqual(test, test_case[i].Literal, statement.Literal)
			assert.DeepEqual(test, test_case[i].Args, statement.Args)
			if len(test_case[i].Children) != len(statement.Children) {
				test.Fatalf(
					"expected to have %d children but got %d",
					len(test_case[i].Children),
					len(statement.Children),
				)
			}
			for j, child := range test_case[i].Children {
				assert.DeepEqual(test, child.Kind, statement.Children[j].Kind)
				assert.DeepEqual(test, child.Literal, statement.Children[j].Literal)
				assert.DeepEqual(test, child.Args, statement.Children[j].Args)
			}
			i++
		}
	})
}
