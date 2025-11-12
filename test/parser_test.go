package libparser_test

import (
	"fmt"
	"path/filepath"
	"testing"

	libparser "github.com/tomefile/lib-parser"
	"gotest.tools/assert"
)

func TestAll(test *testing.T) {
	defer libparser.CloseAll()

	for _, test_case := range ExpectedData {
		path := filepath.Join("data", test_case.Filename)

		test.Run(test_case.Filename, func(test *testing.T) {
			file, err := libparser.OpenFile(path)
			assert.NilError(test, err)

			parser := libparser.New(file).With(libparser.PostNoShebang)

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
	defer libparser.CloseAll()

	message := "this is a result of post processing!"
	path := filepath.Join("data", "01_basic.tome")

	file, err := libparser.OpenFile(path)
	assert.NilError(test, err)

	parser := libparser.New(file).
		With(libparser.PostNoShebang).
		With(libparser.PostExclude[*libparser.ExecNode]).
		With(libparser.PostExclude[*libparser.DirectiveNode]).
		With(
			func(node libparser.Node) (libparser.Node, *libparser.DetailedError) {
				switch node := node.(type) {
				case *libparser.CommentNode:
					node.Contents = message
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
