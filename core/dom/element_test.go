package dom_test

import (
	"testing"

	"github.com/saman3d/samdoc/xml"
	"github.com/saman3d/samtui/core/dom"
	"github.com/stretchr/testify/assert"
)

type elementXMLUnmarshalTestSuite struct {
	name     string
	template string
	expected *dom.Element
}

var elementUnmarshalTestSuites = []elementXMLUnmarshalTestSuite{
	{
		name: "simple_unmarshal",
		template: `
			<element>
				<name>test</name>
				<type>text</type>
				<value>test</value>
			</element>
		`,
		expected: &dom.Element{
			Name:  "element",
			Attrs: dom.NewAttributes(dom.WithDefaultAttributes()),
			Children: []*dom.Element{
				{
					Name:    "name",
					Content: "test",
					Attrs:   dom.NewAttributes(dom.WithDefaultAttributes()),
				},
				{
					Name:    "type",
					Content: "text",
					Attrs:   dom.NewAttributes(dom.WithDefaultAttributes()),
				},
				{
					Name:    "value",
					Content: "test",
					Attrs:   dom.NewAttributes(dom.WithDefaultAttributes()),
				},
			},
		},
	},
	{
		name: "simple_unmarshal_with_attributes",
		template: `
			<element id="test">
				<name id="test_name">test</name>
				<type id="test_type">text</type>
				<value id="test_value">text</value>
		    </element>
		`,
		expected: &dom.Element{
			Name: "element",
			Attrs: &dom.Attributes{
				ID: "test",
			},
			Children: []*dom.Element{
				{
					Name: "name",
					Attrs: &dom.Attributes{
						ID: "test_name",
					},
					Content: "test",
				},
				{
					Name: "type",
					Attrs: &dom.Attributes{
						ID: "test_type",
					},
					Content: "text",
				},
				{
					Name: "value",
					Attrs: &dom.Attributes{
						ID: "test_value",
					},
					Content: "text",
				},
			},
		},
	},
}

func TestElementXMLUnmarshal(t *testing.T) {
	for _, test := range elementUnmarshalTestSuites {
		t.Run(test.name, func(t *testing.T) {
			e := &dom.Element{}
			err := xml.Unmarshal([]byte(test.template), e)
			assert.Nil(t, err)
			assert.Equal(t, test.expected.Name, e.Name)
			for i, e := range e.Children {
				assert.Equal(t, test.expected.Children[i].Name, e.Name)
				assert.Equal(t, test.expected.Children[i].Content, e.Content)
				assert.Equal(t, test.expected.Children[i].Attrs, e.Attrs)
			}
		})
	}
}

func TestElementFromString(t *testing.T) {
	e, err := dom.NewElementFromString(`<element color="2"></element>`)
	assert.Nil(t, err)
	assert.Equal(t, "element", e.Name)
	assert.Equal(t, 2, e.Attrs.Color)
}
