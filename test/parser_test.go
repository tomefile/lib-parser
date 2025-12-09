package libparser_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	libescapes "github.com/bbfh-dev/lib-ansi-escapes"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	liberrors "github.com/tomefile/lib-errors"
	libparser "github.com/tomefile/lib-parser"
	"gotest.tools/assert"
)

var IgnoredOptions = []cmp.Option{
	cmpopts.IgnoreTypes(libparser.NodeContext{}),
	cmpopts.IgnoreFields(libparser.StringModifier{}, "Call"),
}

func TestAll(test *testing.T) {
	defer libparser.CloseAll()

	for _, test_case := range ExpectedData {
		path := filepath.Join("data", test_case.Filename)

		test.Run(test_case.Filename, func(test *testing.T) {
			file, err := libparser.OpenFile(path)
			assert.NilError(test, err)

			parser := libparser.New(file)
			parser.Hooks = []libparser.Hook{libparser.NoShebangHook}

			derr := parser.Run()
			if os.Getenv("GO_DEBUG") == "1" {
				data, _ := json.MarshalIndent(parser.Result.NodeChildren, "", "  ")
				fmt.Fprintf(
					test.Output(),
					"%s%s%s\n",
					libescapes.TextColorWhite,
					string(data),
					libescapes.ColorReset,
				)
			}
			if derr != nil {
				derr.Print(test.Output())
				test.FailNow()
			}
			tree := parser.Result

			assert.DeepEqual(
				test,
				test_case.Expect.NodeChildren,
				tree.NodeChildren,
				IgnoredOptions...)
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

	parser := libparser.New(file)
	parser.Hooks = []libparser.Hook{
		libparser.NoShebangHook,
		libparser.ExcludeHook[*libparser.NodeExec],
		libparser.ExcludeHook[*libparser.NodeDirective],
		func(node libparser.Node) (libparser.Node, *liberrors.DetailedError) {
			switch node := node.(type) {
			case *libparser.NodeComment:
				node.Contents = message
			}
			return node, nil
		},
	}

	if derr := parser.Run(); derr != nil {
		derr.Print(test.Output())
		test.FailNow()
	}
	tree := parser.Result

	assert.DeepEqual(test, tree.NodeChildren, libparser.NodeChildren{
		&libparser.NodeComment{
			Contents: message,
		},
	}, IgnoredOptions...)
}
