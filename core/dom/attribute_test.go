package dom

import (
	"testing"

	"gotest.tools/v3/assert"
)

type attributeParseTestSuite struct {
	name     string
	input    RawAttributeList
	expected *Attributes
}

var attributeParseTestSuites = []attributeParseTestSuite{
	{
		name:     "empty",
		input:    RawAttributeList{},
		expected: &Attributes{},
	},
	{
		name: "single",
		input: RawAttributeList{
			{
				"id",
				"saman",
			},
		},
		expected: &Attributes{
			ID: "saman",
		},
	},
	{
		name: "multiple",
		input: RawAttributeList{
			{
				"id",
				"saman",
			},
			{
				"display",
				"block",
			},
		},
		expected: &Attributes{
			ID:      "saman",
			Display: Display_Block,
		},
	},
	{
		name: "multiple with invalid",
		input: RawAttributeList{
			{
				"id",
				"saman",
			},
			{
				"display",
				"block",
			},
			{
				"position",
				"invalid",
			},
		},
		expected: &Attributes{
			ID:      "saman",
			Display: Display_Block,
		},
	},
}

func TestParse(t *testing.T) {
	for _, suite := range attributeParseTestSuites {
		t.Run(suite.name, func(t *testing.T) {
			actual := NewAttributes()
			actual.Parse(suite.input)
			assert.DeepEqual(t, suite.expected, actual)
		})
	}
}
