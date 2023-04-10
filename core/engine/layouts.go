package engine

import (
	"context"

	"github.com/saman3d/samtui/core/dom"
)

type Layout interface {
	Layout(ctx context.Context, elem *dom.Element, boundry dom.Boundry) error
}

type Flex struct {
	view    View
	rndstck RenderStack
}

func newFlexLayout(v View, rndstck RenderStack) *Flex {
	return &Flex{
		view:    v,
		rndstck: rndstck,
	}
}

func (f *Flex) Layout(ctx context.Context, elem *dom.Element, boundry dom.Boundry) error {
	renderBase(elem, f.view)
	drawBorder(elem, f.view)
	switch elem.Attrs.FlexDirection {
	case dom.FlexDirection_Row:
		return f.renderRow(ctx, elem, boundry)
	case dom.FlexDirection_Column:
		return f.renderColumn(ctx, elem, boundry)
	default:
		return f.renderRow(ctx, elem, boundry)
	}
}

func (f *Flex) renderRow(ctx context.Context, elem *dom.Element, boundry dom.Boundry) error {
	var total_width = boundry.Width()
	var flags = newRenderMan(len(elem.Children), RenderFlag_Flex)
	var total_flex_shares int
	for i, child := range elem.Children {
		if child.Attrs.Display == dom.Display_Absolute {
			flags.SetFlag(i, RenderFlag_Ignore)
			continue
		}
		if child.Attrs.Width != 0 {
			flags.SetFlag(i, RenderFlag_Width|RenderFlag_Calculated)
			flags.SaveCalculation(i, float32(child.Attrs.Width))
			total_width -= child.Attrs.Width
			continue
		} else if child.Attrs.Flex != 0 {
			total_flex_shares += child.Attrs.Flex
			if child.Attrs.MinWidth != 0 || child.Attrs.MaxWidth != 0 {
				flags.SetFlag(i, RenderFlag_MinWidth|RenderFlag_MaxWidth)
				continue
			}
		}
	}

	flex_unit := float32(total_width) / float32(total_flex_shares)

	for head, _, ok := flags.Next(-1, RenderFlag_MaxWidth|RenderFlag_MinWidth); ok; head, _, ok = flags.Next(head, RenderFlag_MaxWidth|RenderFlag_MinWidth) {
		flex_width := flex_unit * float32(elem.Children[head].Attrs.Flex)
		if elem.Children[head].Attrs.MinWidth > 0 && elem.Children[head].Attrs.MinWidth > int(flex_width) {
			flags.SaveCalculation(head, float32(elem.Children[head].Attrs.MinWidth))
			flags.SetFlag(head, RenderFlag_Calculated)
			total_width -= elem.Children[head].Attrs.MinWidth
			total_flex_shares -= elem.Children[head].Attrs.Flex
		} else if elem.Children[head].Attrs.MaxWidth >= 0 && elem.Children[head].Attrs.MaxWidth < int(flex_width) {
			flags.SaveCalculation(head, float32(elem.Children[head].Attrs.MaxWidth))
			flags.SetFlag(head, RenderFlag_Calculated)
			total_width -= elem.Children[head].Attrs.MaxWidth
			total_flex_shares -= elem.Children[head].Attrs.Flex
		}
	}

	reserve := 0
	cap := 0
	for head, d, ok := flags.Next(-1, 0); ok; head, d, ok = flags.Next(head, 0) {
		if d.Flag&RenderFlag_Ignore > 1 {
			f.rndstck.Push(elem.Children[head])
			continue
		}
		if (RenderFlag_Calculated & d.Flag) > 0 {
			elem.Children[head].Boundry = dom.NewBoundry(
				reserve+boundry.FirstX,
				boundry.FirstY,
				reserve+boundry.FirstX+int(d.CalculatedValue),
				boundry.SecondY,
			)
			f.rndstck.Push(elem.Children[head])
			reserve += int(d.CalculatedValue)
			continue
		}
		flex_unit = float32(total_width-cap) / float32(total_flex_shares)
		elem.Children[head].Boundry = dom.NewBoundry(
			reserve+boundry.FirstX,
			boundry.FirstY,
			reserve+boundry.FirstX+int(flex_unit*float32(elem.Children[head].Attrs.Flex)),
			boundry.SecondY,
		)
		f.rndstck.Push(elem.Children[head])
		total_flex_shares -= elem.Children[head].Attrs.Flex
		reserve += int(flex_unit * float32(elem.Children[head].Attrs.Flex))
		cap += int(flex_unit * float32(elem.Children[head].Attrs.Flex))
	}

	return nil
}

func (f *Flex) renderColumn(ctx context.Context, elem *dom.Element, boundry dom.Boundry) error {
	var total_height = boundry.Height()
	var flags = newRenderMan(len(elem.Children), RenderFlag_Flex)
	var total_flex_shares int
	for i, child := range elem.Children {
		if child.Attrs.Display == dom.Display_Absolute {
			flags.SetFlag(i, RenderFlag_Ignore)
			continue
		}
		if child.Attrs.Height != 0 {
			flags.SetFlag(i, RenderFlag_Height|RenderFlag_Calculated)
			flags.SaveCalculation(i, float32(child.Attrs.Height))
			total_height -= child.Attrs.Height
		} else if child.Attrs.Flex != 0 {
			total_flex_shares += child.Attrs.Flex
			if child.Attrs.MinHeight != 0 || child.Attrs.MaxHeight != 0 {
				flags.SetFlag(i, RenderFlag_MinHeight|RenderFlag_MaxHeight)
				continue
			}
		}
	}

	flex_unit := float32(total_height) / float32(total_flex_shares)
	for head, _, ok := flags.Next(-1, RenderFlag_MaxHeight|RenderFlag_MinHeight); ok; head, _, ok = flags.Next(head, RenderFlag_MaxHeight|RenderFlag_MinHeight) {
		flex_height := flex_unit * float32(elem.Children[head].Attrs.Flex)
		if elem.Children[head].Attrs.MinHeight > 0 && elem.Children[head].Attrs.MinHeight > int(flex_height) {
			flags.SaveCalculation(head, float32(elem.Children[head].Attrs.MinHeight))
			flags.SetFlag(head, RenderFlag_Calculated)
			total_height -= elem.Children[head].Attrs.MinHeight
			total_flex_shares -= elem.Children[head].Attrs.Flex
		} else if elem.Children[head].Attrs.MaxHeight > 0 && elem.Children[head].Attrs.MaxHeight <= int(flex_height) {
			flags.SaveCalculation(head, float32(elem.Children[head].Attrs.MaxHeight))
			flags.SetFlag(head, RenderFlag_Calculated)
			total_height -= elem.Children[head].Attrs.MaxHeight
			total_flex_shares -= elem.Children[head].Attrs.Flex
		}
	}

	reserve := 0
	cap := 0
	for head, d, ok := flags.Next(-1, 0); ok; head, d, ok = flags.Next(head, 0) {
		if d.Flag&RenderFlag_Ignore > 1 {
			f.rndstck.Push(elem.Children[head])
			continue
		}
		if head < elem.State.ScrollY {
			continue
		}
		if (RenderFlag_Calculated & d.Flag) > 0 {
			if reserve+int(d.CalculatedValue) > boundry.Height() {
				break
			}
			elem.Children[head].Boundry = dom.NewBoundry(
				boundry.FirstX,
				reserve+boundry.FirstY,
				boundry.SecondX,
				reserve+boundry.FirstY+int(d.CalculatedValue),
			)
			reserve += int(d.CalculatedValue)
			goto rend
		}
		if reserve+int(flex_unit*float32(elem.Children[head].Attrs.Flex)) > boundry.Height() {
			break
		}
		flex_unit = float32(total_height-cap) / float32(total_flex_shares)
		elem.Children[head].Boundry = dom.NewBoundry(
			boundry.FirstX,
			reserve+boundry.FirstY,
			boundry.SecondX,
			reserve+boundry.FirstY+int(flex_unit*float32(elem.Children[head].Attrs.Flex)),
		)
		total_flex_shares -= elem.Children[head].Attrs.Flex
		reserve += int(flex_unit * float32(elem.Children[head].Attrs.Flex))
		cap += int(flex_unit * float32(elem.Children[head].Attrs.Flex))

	rend:
		f.rndstck.Push(elem.Children[head])
	}

	return nil
}

type Block struct {
	View    View
	rndstck RenderStack
}

func newBlockLayout(v View, rndstck RenderStack) *Block {
	return &Block{
		View:    v,
		rndstck: rndstck,
	}
}

func (b *Block) Layout(ctx context.Context, elem *dom.Element, boundry dom.Boundry) error {
	renderBase(elem, b.View)
	drawBorder(elem, b.View)
	renderText(elem, b.View)
	// for _, child := range elem.Children {
	// 	b.renderchan <- newBoundedElement(child, boundry.Shrink(1))
	// }
	return nil
}

type Absolute struct {
	View    View
	rndstck RenderStack
}

func newAbsoluteLayout(v View, rndstck RenderStack) *Absolute {
	return &Absolute{
		View:    v,
		rndstck: rndstck,
	}
}

func (a *Absolute) Layout(ctx context.Context, elem *dom.Element, boundry dom.Boundry) error {
	elem.Boundry = dom.NewBoundry(elem.Attrs.Left, elem.Attrs.Top, elem.Attrs.Width+elem.Attrs.Left, elem.Attrs.Height+elem.Attrs.Top)
	renderBase(elem, a.View)
	drawBorder(elem, a.View)
	renderText(elem, a.View)
	return nil
}

type LayoutType string

const (
	LayoutType_Flex     LayoutType = "flex"
	LayoutType_Block    LayoutType = "block"
	LayoutType_Absolute LayoutType = "absolute"
)

type RenderFlag uint16

const (
	RenderFlag_Flex RenderFlag = 1 << iota
	RenderFlag_Width
	RenderFlag_MaxWidth
	RenderFlag_MinWidth
	RenderFlag_Height
	RenderFlag_MaxHeight
	RenderFlag_MinHeight
	RenderFlag_Calculated
	RenderFlag_Ignore
	RenderFlag_Invalid
)

type RenderData struct {
	Flag            RenderFlag
	CalculatedValue float32
}

func newRednerData(defaultFlag RenderFlag) *RenderData {
	return &RenderData{
		Flag: defaultFlag,
	}
}

type RenderMan map[int]*RenderData

func newRenderMan(l int, defaultFlag RenderFlag) RenderMan {
	res := make(RenderMan, 0)
	for x := 0; x < l; x++ {
		res[x] = newRednerData(defaultFlag)
	}
	return res
}

func (rfl RenderMan) SetFlag(i int, flag RenderFlag) {
	rfl[i].Flag = flag
}

func (rfl RenderMan) Flag(i int) RenderFlag {
	return rfl[i].Flag
}

func (rfl RenderMan) SaveCalculation(i int, value float32) {
	rfl[i].CalculatedValue = value
}

func (rfl RenderMan) Next(head int, flag RenderFlag) (int, *RenderData, bool) {
	head++
	for d, ok := rfl[head]; ok; d, ok = rfl[head] {
		if flag == 0 || (flag&d.Flag) > 0 {
			return head, d, ok
		}
		head++
	}
	return head, nil, false
}

func (rfl RenderMan) FlagRest(head int, flag RenderFlag) {
	for d, ok := rfl[head]; ok; d, ok = rfl[head] {
		d.Flag = flag
		head++
	}
}
