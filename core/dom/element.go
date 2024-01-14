package dom

import (
	"fmt"

	"github.com/saman3d/samdoc/xml"
	"github.com/saman3d/samtui/core/common"
)

type Element struct {
	Name     string
	Children []*Element
	Parent   *Element
	Style    []rune
	Content  string
	Attrs    *Attributes
	Boundry  Boundry
	State    *ElementState
}

func NewElement(name string) *Element {
	return &Element{
		Name:     name,
		Children: make([]*Element, 0),
		Attrs:    NewAttributes(),
		Boundry:  NewBoundry(-1, -1, -1, -1),
		State:    NewElementState(),
	}
}

func NewElementFromString(s string) (*Element, error) {
	el := NewElement("root")
	err := xml.Unmarshal([]byte(s), el)
	if err != nil {
		return nil, err
	}
	return el, nil
}

func MustParseElementFromString(s string) *Element {
	el := NewElement("root")
	err := xml.Unmarshal([]byte(s), el)
	if err != nil {
		panic(err)
	}
	return el
}

func (el *Element) XMLUnmarshal(d *xml.XMLDecoder, start xml.StartTag) error {
	el.Name = start.Tagname
	el.Attrs = NewAttributes(WithDefaultAttributes())
	if el.Parent != nil {
		el.Attrs.InheritFrom(el.Parent.Attrs)
	}
	el.Attrs.Parse(rawAttrsToAttibuteList(start.Attrs))
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch tok := tok.(type) {
		case xml.CharData:
			el.Content += string(tok)

		case xml.StartTag:
			var e = NewElement(tok.Tagname)
			e.Parent = el
			err = e.XMLUnmarshal(d, tok)
			if err != nil {
				return err
			}
			el.Children = append(el.Children, e)

		case xml.EndTag:
			if tok.Tagname != start.Tagname {
				return common.ErrUnexpectedEndTag
			}
			return nil
		}
	}
}

func (el *Element) AppendChild(e *Element) {
	el.Children = append(el.Children, e)
}

func (el *Element) InheritChildrensAttr() {
	for _, c := range el.Children {
		c.Attrs.InheritFrom(el.Attrs)
		c.InheritChildrensAttr()
	}
}

type Boundry struct {
	FirstX  int
	FirstY  int
	SecondX int
	SecondY int
}

func NewBoundry(fx, fy, sx, sy int) Boundry {
	return Boundry{
		FirstX:  fx,
		FirstY:  fy,
		SecondX: sx,
		SecondY: sy,
	}
}

func (b Boundry) Indexify() Boundry {
	b.SecondX -= 1
	b.SecondY -= 1
	return b
}

func (b Boundry) Normalize() Boundry {
	if b.FirstX <= 0 {
		b.FirstX = 1
	}
	if b.FirstY <= 0 {
		b.FirstY = 1
	}
	return b
}

func (b Boundry) Shrink(by int) Boundry {
	b.FirstX += by
	b.FirstY += by
	b.SecondX -= by
	b.SecondY -= by
	return b
}

func (b Boundry) ShrinkMask(by int, m PositionMask) Boundry {
	if m&PositionMaskTop != 0 {
		b.FirstY += by
	}
	if m&PositionMaskBottom != 0 {
		b.SecondY -= by
	}
	if m&PositionMaskLeft != 0 {
		b.FirstX += by
	}
	if m&PositionMaskRight != 0 {
		b.SecondX -= by
	}
	return b
}

type PositionMask byte

const (
	PositionMaskTop PositionMask = 1 << iota
	PositionMaskBottom
	PositionMaskLeft
	PositionMaskRight
)

func (b Boundry) Inflate(by int) Boundry {
	b.FirstX -= by
	b.FirstY -= by
	b.SecondX += by
	b.SecondY += by
	return b
}

func (b Boundry) InflateMask(by int, m PositionMask) Boundry {
	if m&PositionMaskTop != 0 {
		b.FirstY -= by
	}
	if m&PositionMaskBottom != 0 {
		b.SecondY += by
	}
	if m&PositionMaskLeft != 0 {
		b.FirstX -= by
	}
	if m&PositionMaskRight != 0 {
		b.SecondX += by
	}
	return b
}

func (b Boundry) Width() int {
	return b.SecondX - b.FirstX
}

func (b Boundry) Height() int {
	return b.SecondY - b.FirstY
}

func (b Boundry) String() string {
	return fmt.Sprintf("Boundry: (%d,%d) (%d,%d)", b.FirstX, b.FirstY, b.SecondX, b.SecondY)
}

func (b Boundry) Circumscribes(b2 Boundry) bool {
	return b.FirstX < b2.FirstX &&
		b.FirstY < b2.FirstY &&
		b.SecondX > b2.SecondX &&
		b.SecondY > b2.SecondY
}

func (b Boundry) Sum(b2 Boundry) Boundry {
	return Boundry{
		FirstX:  min(b.FirstX, b2.FirstX),
		FirstY:  min(b.FirstY, b2.FirstY),
		SecondX: max(b.SecondX, b2.SecondX),
		SecondY: max(b.SecondY, b2.SecondY),
	}
}

type ElementState struct {
	ScrollX int
	ScrollY int
}

func NewElementState() *ElementState {
	return &ElementState{}
}

func (s *ElementState) ScrollTo(x, y int) {
	s.ScrollX = x
	s.ScrollY = y
}

func (s *ElementState) ScrollBy(x, y int) {
	s.ScrollX += x
	s.ScrollY += y
}
