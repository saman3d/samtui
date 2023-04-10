package engine

import (
	"github.com/saman3d/samtui/core/dom"
	"github.com/saman3d/samtui/core/engine/view"
)

func drawBorder(elem *dom.Element, v View) {
	if !elem.Attrs.Border {
		return
	}
	boundry := elem.Boundry.Indexify()
	v.PrintRuneRepeat(boundry.FirstX, boundry.FirstY, elem.Attrs.Color, elem.Attrs.BackGroundColor, boundry.Width(), elem.Attrs.ZIndex, view.AxisMask_X, '─')
	v.PrintRuneRepeat(boundry.FirstX, boundry.SecondY, elem.Attrs.Color, elem.Attrs.BackGroundColor, boundry.Width(), elem.Attrs.ZIndex, view.AxisMask_X, '─')
	v.PrintRuneRepeat(boundry.FirstX, boundry.FirstY, elem.Attrs.Color, elem.Attrs.BackGroundColor, boundry.Height(), elem.Attrs.ZIndex, view.AxisMask_Y, '│')
	v.PrintRuneRepeat(boundry.SecondX, boundry.FirstY, elem.Attrs.Color, elem.Attrs.BackGroundColor, boundry.Height(), elem.Attrs.ZIndex, view.AxisMask_Y, '│')

	v.PrintRune(boundry.FirstX, boundry.FirstY, elem.Attrs.Color, elem.Attrs.BackGroundColor, elem.Attrs.ZIndex, '┌')
	v.PrintRune(boundry.SecondX, boundry.FirstY, elem.Attrs.Color, elem.Attrs.BackGroundColor, elem.Attrs.ZIndex, '┐')

	v.PrintRune(boundry.FirstX, boundry.SecondY, elem.Attrs.Color, elem.Attrs.BackGroundColor, elem.Attrs.ZIndex, '└')
	v.PrintRune(boundry.SecondX, boundry.SecondY, elem.Attrs.Color, elem.Attrs.BackGroundColor, elem.Attrs.ZIndex, '┘')
}

func renderText(elem *dom.Element, v View) dom.Boundry {
	boundry := elem.Boundry
	if elem.Attrs.Border {
		boundry = boundry.Shrink(1)
	}
	text := elem.Content
	n := 0
	if text == "" || boundry.Width() < 1 || boundry.Height() < 1 {
		return boundry
	}

	boundry = boundry.Indexify()

c:
	for y := 0; y <= boundry.Height(); y++ {
		if y >= elem.State.ScrollY {
			for x := 0; x <= boundry.Width(); x++ {
				if x >= len(text) {
					break c
				}
				v.PrintRune(boundry.FirstX+x, boundry.FirstY+y, elem.Attrs.Color, elem.Attrs.BackGroundColor, elem.Attrs.ZIndex, rune(text[x]))
			}
			text = text[boundry.Width()+1:]
		}
		n++
	}
	return boundry.ShrinkMask(n, dom.PositionMaskTop)
}

func renderBase(elem *dom.Element, v View) {
	for y := elem.Boundry.FirstY; y < elem.Boundry.SecondY; y++ {
		for x := elem.Boundry.FirstX; x < elem.Boundry.SecondX; x++ {
			v.PrintRune(x, y, elem.Attrs.Color, elem.Attrs.BackGroundColor, elem.Attrs.ZIndex, ' ')
		}
	}
}
