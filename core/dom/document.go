package dom

import (
	"errors"
	"io"

	"github.com/saman3d/samdoc/xml"
)

var (
	ErrorUnexpectedEndTag = errors.New("unexpected end tag")
)

type Document struct {
	Head *Head
	Body *Element `xml:",any"`
}

func NewDocumentFromReader(f io.Reader) (*Document, error) {
	doc := &Document{
		Head: &Head{},
		Body: NewElement("body"),
	}

	fb, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(fb, doc); err != nil && err != io.EOF {
		return nil, err
	}
	return doc, nil
}

func (doc *Document) XMLUnmarshal(d *xml.XMLDecoder, start xml.StartTag) error {
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch tok := tok.(type) {
		case xml.StartTag:
			switch tok.Tagname {
			case "head":
				doc.Head = &Head{}
				if err := doc.Head.XMLUnmarshal(d, tok); err != nil {
					return err
				}
			case "body":
				doc.Body = NewElement("body")
				if err := doc.Body.XMLUnmarshal(d, tok); err != nil {
					return err
				}
			}
		case xml.EndTag:
			if tok.Tagname != "html" {
				return ErrorUnexpectedEndTag
			}
			return nil
		}
	}
}

type Head struct {
	Title string
}

func (head *Head) XMLUnmarshal(d *xml.XMLDecoder, start xml.StartTag) error {
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch tok := tok.(type) {
		case xml.StartTag:
			switch tok.Tagname {
			case "title":
				var unel xml.UniversalElement
				err := unel.XMLUnmarshal(d, tok)
				if err != nil {
					return err
				}
				head.Title = unel.Data
			}
		case xml.EndTag:
			return nil
		}
	}
}
