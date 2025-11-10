package libparser_test

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	libparser "github.com/tomefile/lib-parser"
	"gotest.tools/assert"
)

func TestAll(test *testing.T) {
	for _, test_case := range ExpectedData {
		file, err := os.OpenFile(
			filepath.Join("data", test_case.Filename),
			os.O_RDONLY,
			os.ModePerm,
		)
		assert.NilError(test, err)
		defer file.Close()

		test.Run(file.Name(), func(test *testing.T) {
			parser := libparser.New(file.Name(), bufio.NewReader(file))
			tree, parser_err := parser.Parse()
			if parser_err != nil {
				fmt.Println(parser_err.GetBeautyPrinted())
				test.FailNow()
			}

			assert.DeepEqual(test, test_case.Expect, tree)
		})
	}
}

func TestPostProcessor(test *testing.T) {
	message := "this is a result of post processing!"

	file, err := os.OpenFile(
		filepath.Join("data", "01_basic.tome"),
		os.O_RDONLY,
		os.ModePerm,
	)
	assert.NilError(test, err)
	defer file.Close()

	parser := libparser.New(
		file.Name(),
		bufio.NewReader(file),
		func(node libparser.Node) (libparser.Node, *libparser.DetailedError) {
			switch node := node.(type) {
			case *libparser.CommentNode:
				node.Contents = message
			}
			return node, nil
		},
		func(node libparser.Node) (libparser.Node, *libparser.DetailedError) {
			switch node.(type) {
			case *libparser.ExecNode, *libparser.DirectiveNode:
				// discard
				return nil, nil
			}
			return node, nil
		},
	)
	tree, parser_err := parser.Parse()
	if parser_err != nil {
		fmt.Println(parser_err.GetBeautyPrinted())
		test.FailNow()
	}

	assert.DeepEqual(test, tree.NodeChildren, libparser.NodeChildren{
		&libparser.CommentNode{
			Contents: message,
		},
	})
}
