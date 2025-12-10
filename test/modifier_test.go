package libparser_test

import (
	"os"
	"testing"

	libparser "github.com/tomefile/lib-parser"
	"gotest.tools/assert"
)

type ModifierTestCase struct {
	Name   libparser.ModifierName
	Input  string
	Args   []*libparser.NodeString
	Expect string
}

var ModifierExpectedData = []ModifierTestCase{
	{
		Name:   libparser.MOD_NOT,
		Input:  "",
		Args:   []*libparser.NodeString{},
		Expect: "1",
	},
	{
		Name:   libparser.MOD_TO_LOWER,
		Input:  "Hello World!",
		Args:   []*libparser.NodeString{},
		Expect: "hello world!",
	},
	{
		Name:   libparser.MOD_TO_UPPER,
		Input:  "Hello World!",
		Args:   []*libparser.NodeString{},
		Expect: "HELLO WORLD!",
	},
	{
		Name:   libparser.MOD_TO_SNAKE,
		Input:  "Hello World!",
		Args:   []*libparser.NodeString{},
		Expect: "hello_world!",
	},
	{
		Name:   libparser.MOD_TO_KEBAB,
		Input:  "Hello World!",
		Args:   []*libparser.NodeString{},
		Expect: "hello-world!",
	},
	{
		Name:   libparser.MOD_TO_CAMEL,
		Input:  "Hello World!",
		Args:   []*libparser.NodeString{},
		Expect: "helloWorld",
	},
	{
		Name:   libparser.MOD_TO_PASCAL,
		Input:  "Hello World!",
		Args:   []*libparser.NodeString{},
		Expect: "HelloWorld",
	},
	{
		Name:   libparser.MOD_TO_DELIMITED,
		Input:  "Hello World!",
		Args:   []*libparser.NodeString{libparser.NewSimpleNodeString(".")},
		Expect: "hello.world!",
	},
	{
		Name:   libparser.MOD_DOES_EXIST,
		Input:  getRealDir(),
		Args:   []*libparser.NodeString{},
		Expect: "1",
	},
	{
		Name:   libparser.MOD_IS_EMPTY,
		Input:  "",
		Args:   []*libparser.NodeString{},
		Expect: "1",
	},
	{
		Name:   libparser.MOD_IS_FILE,
		Input:  getRealDir(),
		Args:   []*libparser.NodeString{},
		Expect: "0",
	},
	{
		Name:   libparser.MOD_IS_DIR,
		Input:  getRealDir(),
		Args:   []*libparser.NodeString{},
		Expect: "1",
	},
	{
		Name:   libparser.MOD_IS_SYMLINK,
		Input:  getRealDir(),
		Args:   []*libparser.NodeString{},
		Expect: "0",
	},
	{
		Name:   libparser.MOD_LENGTH,
		Input:  "abcdef",
		Args:   []*libparser.NodeString{},
		Expect: "6",
	},
	{
		Name:   libparser.MOD_QUOTED,
		Input:  "abcdef",
		Args:   []*libparser.NodeString{},
		Expect: "\"abcdef\"",
	},
	{
		Name:   libparser.MOD_TRIM,
		Input:  "++abc+",
		Args:   []*libparser.NodeString{libparser.NewSimpleNodeString("+")},
		Expect: "abc",
	},
	{
		Name:   libparser.MOD_TRIM_PREFIX,
		Input:  "++abc+",
		Args:   []*libparser.NodeString{libparser.NewSimpleNodeString("+")},
		Expect: "+abc+",
	},
	{
		Name:   libparser.MOD_TRIM_SUFFIX,
		Input:  "++abc+",
		Args:   []*libparser.NodeString{libparser.NewSimpleNodeString("+")},
		Expect: "++abc",
	},
	{
		Name:  libparser.MOD_PAD,
		Input: "abc",
		Args: []*libparser.NodeString{
			libparser.NewSimpleNodeString("3"),
			libparser.NewSimpleNodeString("5"),
		},
		Expect: "   abc     ",
	},
	{
		Name:   libparser.MOD_PAD_LEFT,
		Input:  "abc",
		Args:   []*libparser.NodeString{libparser.NewSimpleNodeString("3")},
		Expect: "   abc",
	},
	{
		Name:   libparser.MOD_PAD_RIGHT,
		Input:  "abc",
		Args:   []*libparser.NodeString{libparser.NewSimpleNodeString("3")},
		Expect: "abc   ",
	},
	{
		Name:   libparser.MOD_HAS_PREFIX,
		Input:  "abcdef",
		Args:   []*libparser.NodeString{libparser.NewSimpleNodeString("abc")},
		Expect: "1",
	},
	{
		Name:   libparser.MOD_HAS_SUFFIX,
		Input:  "abcdef",
		Args:   []*libparser.NodeString{libparser.NewSimpleNodeString("abc")},
		Expect: "0",
	},
	{
		Name:  libparser.MOD_SLICE,
		Input: "abcdef",
		Args: []*libparser.NodeString{
			libparser.NewSimpleNodeString("1"),
			libparser.NewSimpleNodeString("-2"),
		},
		Expect: "bcd",
	},
	{
		Name:   libparser.MOD_REVERSE,
		Input:  "abc",
		Args:   []*libparser.NodeString{},
		Expect: "cba",
	},
	{
		Name:   libparser.MOD_INVERT,
		Input:  "Hello World!",
		Args:   []*libparser.NodeString{},
		Expect: "hELLO wORLD!",
	},
}

func TestModifier(test *testing.T) {
	for _, test_case := range ModifierExpectedData {
		test.Run(string(test_case.Name), func(test *testing.T) {
			modifier, err := libparser.GetModifier(test_case.Name, test_case.Args)
			assert.NilError(test, err)
			assert.Assert(test, modifier.Call != nil)
			assert.DeepEqual(
				test,
				test_case.Expect,
				modifier.Call(libparser.Locals{}, test_case.Input),
			)
		})
	}
}

func getRealDir() string {
	dir, _ := os.Getwd()
	return dir
}
