package libparser_test

import (
	"encoding/json"
	"testing"

	libparser "github.com/tomefile/lib-parser"
	"gotest.tools/assert"
)

func TestMarshalJSON(test *testing.T) {
	modifier := libparser.GetModifier("to_lower", []string{})
	data, err := json.Marshal(modifier)
	assert.NilError(test, err)
	assert.DeepEqual(test, string(data), `"modifier()"`)
}
