package dom_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/saman3d/samtui/core/dom"
	"gotest.tools/v3/assert"
)

type domFromReaderTestSuite struct {
	name     string
	input    string
	expected *dom.Document
}

var domFromReaderTestSuites = []domFromReaderTestSuite{
	{
		name: "simple",
		input: `
<!DOCTYPE html>
<html>
<head>
	<title>Test</title>
</head>
<body>
	<div id="test">
		<div id="test2"></div>
	</div>
</body>
</html>
`,
		expected: &dom.Document{
			Head: &dom.Head{
				Title: "Test",
			},
			Body: &dom.Element{
				Name:  "body",
				Attrs: dom.NewAttributes(),
				Children: []*dom.Element{
					{
						Name: "div",
						Attrs: &dom.Attributes{
							ID: "test",
						},
						Children: []*dom.Element{
							{
								Name: "div",
								Attrs: &dom.Attributes{
									ID: "test2",
								},
							},
						},
					},
				},
			},
		},
	},
}

func TestDOMFromReader(t *testing.T) {
	for _, test := range domFromReaderTestSuites {
		t.Run(test.name, func(t *testing.T) {
			r := strings.NewReader(test.input)
			doc, err := dom.NewDocumentFromReader(r)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			assert.DeepEqual(t, doc, test.expected, ignoreParent())
		})
	}
}

func ignoreParent() cmp.Option {
	return cmp.FilterPath(func(p cmp.Path) bool {
		return p.Last().String() == ".Parent"
	}, cmp.Ignore())
}
