package parser_test

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	libparser "github.com/tomefile/lib-parser"
	"gotest.tools/assert"
)

var NoChildren = []*libparser.Node{}

var FileTestCases = map[string][]*libparser.Node{
	"01_basic.tome": {
		{
			Type:    libparser.NODE_COMMENT,
			Literal: " Example program, привет мир 🟡!",
		},
		{
			Type:     libparser.NODE_DIRECTIVE,
			Literal:  "include",
			Args:     []any{"<std>"},
			Children: NoChildren,
		},
		{
			Type:    libparser.NODE_EXEC,
			Literal: "echo",
			Args:    []any{"Hello World!", "and another line", "and another."},
		},
	},
	"02_directive_body.tome": {
		{
			Type:    libparser.NODE_EXEC,
			Literal: "echo",
			Args:    []any{"1"},
		},
		{
			Type:    libparser.NODE_DIRECTIVE,
			Literal: "section",
			Args:    []any{"Hello World!"},
			Children: []*libparser.Node{
				{
					Type:    libparser.NODE_EXEC,
					Literal: "echo",
					Args:    []any{"1.1"},
				},
				{
					Type:    libparser.NODE_EXEC,
					Literal: "echo",
					Args:    []any{"1.2"},
				},
			},
		},
	},
	"03_directive_nested.tome": {
		{
			Type:    libparser.NODE_EXEC,
			Literal: "echo",
			Args:    []any{"1"},
		},
		{
			Type:    libparser.NODE_DIRECTIVE,
			Literal: "section",
			Args:    []any{"Hello World!"},
			Children: []*libparser.Node{
				{
					Type:    libparser.NODE_EXEC,
					Literal: "echo",
					Args:    []any{"1.1"},
				},
				{
					Type:    libparser.NODE_EXEC,
					Literal: "echo",
					Args:    []any{"1.2"},
				},
				{
					Type:    libparser.NODE_DIRECTIVE,
					Literal: "section",
					Args:    []any{"Nested"},
					Children: []*libparser.Node{
						{
							Type:    libparser.NODE_COMMENT,
							Literal: " This is nested inside",
						},
						{
							Type:    libparser.NODE_EXEC,
							Literal: "echo",
							Args:    []any{"2.1"},
						},
						{
							Type:    libparser.NODE_EXEC,
							Literal: "echo",
							Args:    []any{"2.2"},
						},
					},
				},
				{
					Type:    libparser.NODE_EXEC,
					Literal: "echo",
					Args:    []any{"1.3"},
				},
			},
		},
	},
	"04_subcommand.tome": {
		{
			Type:    libparser.NODE_EXEC,
			Literal: "echo",
			Args: []any{
				"123",
				&libparser.Node{
					Type:    libparser.NODE_EXEC,
					Literal: "readlink",
					Args:    []any{"-p", "$MY_LINK"},
				},
				"456",
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

		consumer := make(chan *libparser.Node)
		parser := libparser.New(bufio.NewReader(file), consumer)

		go parser.ParseComplete()

		var i int
		for node := range consumer {
			// data, err := json.MarshalIndent(node, "", strings.Repeat(" ", 4))
			// assert.NilError(test, err)
			// test.Logf(
			// 	"%s\nRaw data: %q",
			// 	string(data),
			// 	string(buffer[node.OffsetStart:node.OffsetEnd]),
			// )

			if len(test_case) <= i {
				test.Fatalf("unknown node: %#v", node)
			}

			// Using individual fields because not every field is important
			assert.DeepEqual(test, test_case[i].Type, node.Type)
			assert.DeepEqual(test, test_case[i].Literal, node.Literal)

			if len(test_case[i].Args) != len(node.Args) {
				test.Fatalf(
					"expected to have %d args but got %d",
					len(test_case[i].Args),
					len(node.Args),
				)
			}
			for j, arg := range test_case[i].Args {
				switch arg := arg.(type) {
				case *libparser.Node:
					assert.DeepEqual(test, arg.Type, node.Args[j].(*libparser.Node).Type)
					assert.DeepEqual(test, arg.Literal, node.Args[j].(*libparser.Node).Literal)
					assert.DeepEqual(test, arg.Args, node.Args[j].(*libparser.Node).Args)
				default:
					assert.DeepEqual(test, arg, node.Args[j])
				}
			}

			if len(test_case[i].Children) != len(node.Children) {
				test.Fatalf(
					"expected to have %d children but got %d",
					len(test_case[i].Children),
					len(node.Children),
				)
			}
			for j, child := range test_case[i].Children {
				assert.DeepEqual(test, child.Type, node.Children[j].Type)
				assert.DeepEqual(test, child.Literal, node.Children[j].Literal)
				assert.DeepEqual(test, child.Args, node.Children[j].Args)
			}
			i++
		}
	})
}

func TestArgParser(test *testing.T) {
	parts, err := libparser.ParseArg("hello $message")
	assert.NilError(test, err)
	assert.DeepEqual(test, parts, []libparser.ArgPart{
		{
			Literal:    "hello ",
			Format:     nil,
			IsVariable: false,
		},
		{
			Literal:    "message",
			Format:     nil,
			IsVariable: true,
		},
	})
}
