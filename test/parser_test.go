package libparser_test

import (
	"path/filepath"
	"testing"

	liberrors "github.com/tomefile/lib-errors"
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

			if derr := parser.Parse(); derr != nil {
				derr.Print(test.Output())
				test.FailNow()
			}
			tree := parser.Result()

			assert.DeepEqual(test, test_case.Expect.NodeChildren, tree.NodeChildren)
			for key := range tree.Tomes {
				// We don't care about what it points to, just that it exists.
				tree.Tomes[key] = nil
			}
			assert.DeepEqual(test, test_case.Expect.Tomes, tree.Tomes)
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
			func(node libparser.Node) (libparser.Node, *liberrors.DetailedError) {
				switch node := node.(type) {
				case *libparser.CommentNode:
					node.Contents = message
				}
				return node, nil
			},
		)

	if derr := parser.Parse(); derr != nil {
		derr.Print(test.Output())
		test.FailNow()
	}
	tree := parser.Result()

	assert.DeepEqual(test, tree.NodeChildren, libparser.NodeChildren{
		&libparser.CommentNode{
			Contents: message,
		},
	})
}
