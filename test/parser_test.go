package libparser_test

import (
	"fmt"
	"path/filepath"
	"testing"

	libparser "github.com/tomefile/lib-parser"
	"gotest.tools/assert"
)

func TestAll(test *testing.T) {
	for _, test_case := range ExpectedData {
		path := filepath.Join("data", test_case.Filename)

		test.Run(test_case.Filename, func(test *testing.T) {
			parser, err := libparser.OpenNew(path)
			assert.NilError(test, err)
			defer parser.Close()

			parser.With(libparser.PostNoShebang)

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

	path := filepath.Join("data", "01_basic.tome")

	parser, err := libparser.OpenNew(path)
	assert.NilError(test, err)
	defer parser.Close()

	parser.With(libparser.PostNoShebang)
	parser.With(
		func(node libparser.Node) (libparser.Node, *libparser.DetailedError) {
			switch node := node.(type) {
			case *libparser.CommentNode:
				node.Contents = message
			}
			return node, nil
		},
	)
	parser.With(libparser.PostExclude[*libparser.ExecNode])
	parser.With(libparser.PostExclude[*libparser.DirectiveNode])

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
