package libparser_test

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
				var builder strings.Builder
				parser_err.BeautyPrint(&builder)
				fmt.Println(builder.String())
				test.FailNow()
			}

			assert.DeepEqual(test, test_case.Expect, tree)
		})
	}
}
