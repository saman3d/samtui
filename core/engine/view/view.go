package view

import (
	"fmt"

	"github.com/saman3d/samtui/core/dom"
)

type View []CellList

type CellList []*Cell

func (cl CellList) String() string {
	var s string
	for _, c := range cl {
		s += c.String()
	}
	return s
}

func (cl CellList) Bytes() []byte {
	var b = make([]byte, 0)
	for _, c := range cl {
		b = append(b, c.Bytes()...)
	}
	return b
}

func NewView(width, height int64) *View {
	view := make([]CellList, height)
	for i := range view {
		view[i] = make(CellList, width)
		for j := range view[i] {
			view[i][j] = newCell()
		}
	}

	v := View(view)
	return &v
}

func (v *View) Width() int64 {
	return int64(len((*v)[0]))
}

func (v *View) Height() int64 {
	return int64(len(*v))
}

func (v *View) Boundry() dom.Boundry {
	return dom.NewBoundry(0, 0, int(v.Width()), int(v.Height()))
}

func (v *View) Flush() {
	for _, row := range *v {
		for _, cell := range row {
			cell.Flush()
		}
	}
}

func (v *View) PrintString(x, y, fg, bg int, zindx uint8, s string) {
	for i, r := range s {
		if (*v)[y][x+i].ZIndex <= zindx {
			(*v)[y][x+i].Style = NewStyle(fg, bg)
			(*v)[y][x+i].Content = r
			(*v)[y][x+i].ZIndex = zindx
		}
	}
}

func (v *View) PrintRune(x, y, fg, bg int, zindx uint8, r rune) {
	if (*v)[y][x].ZIndex <= zindx {
		(*v)[y][x].ZIndex = zindx
		(*v)[y][x].Style = NewStyle(fg, bg)
		(*v)[y][x].Content = r
	}
}

func (v *View) PrintRuneRepeat(x, y, fg, bg, rp int, zindx uint8, axis AxisMask, r rune) {
	switch axis {
	case AxisMask_X:
		for i := 0; i < rp; i++ {
			if (*v)[y][x+i].ZIndex <= zindx {
				(*v)[y][x+i].Style = NewStyle(fg, bg)
				(*v)[y][x+i].Content = r
				(*v)[y][x+i].ZIndex = zindx
			}
		}
	case AxisMask_Y:
		for i := 0; i < rp; i++ {
			if (*v)[y+i][x].ZIndex <= zindx {
				(*v)[y+i][x].Style = NewStyle(fg, bg)
				(*v)[y+i][x].Content = r
				(*v)[y+i][x].ZIndex = zindx
			}
		}
	case AxisMask_X | AxisMask_Y:
		for i := 0; i < rp; i++ {
			for j := 0; j < rp; j++ {
				if (*v)[y+i][x+j].ZIndex <= zindx {
					(*v)[y+i][x+j].Style = NewStyle(fg, bg)
					(*v)[y+i][x+j].Content = r
					(*v)[y+i][x+j].ZIndex = zindx
				}
			}
		}
	default:
		for i := 0; i < rp; i++ {
			if (*v)[y][x+i].ZIndex <= zindx {
				(*v)[y][x-i].Style = NewStyle(fg, bg)
				(*v)[y][x+i].Content = r
				(*v)[y][x+i].ZIndex = zindx
			}
		}
	}
}

func (v *View) Slice(x, y, l int) CellList {
	// fmt.Println(x, y, l)
	return (*v)[y][x : x+l]
}

func (v *View) Render() {
	for _, row := range *v {
		for _, cell := range row {
			fmt.Print(string(cell.Content))
		}
	}
}

func (v *View) ClearBoundry(bndr dom.Boundry) {
	for y := bndr.FirstY; y < bndr.SecondY; y++ {
		for x := bndr.FirstX; x < bndr.SecondX; x++ {
			(*v)[y][x].Content = ' '
			(*v)[y][x].Style = NewStyle(0, 0)
			(*v)[y][x].ZIndex = 0
		}
	}
}

func (v *View) FillBoundry(fg, bg int, bndr dom.Boundry) {
	s := NewStyle(fg, bg)
	for y := bndr.FirstY; y < bndr.SecondY; y++ {
		for x := bndr.FirstX; x < bndr.SecondX; x++ {
			(*v)[y][x].Style = s
		}
	}
}

func (v *View) Resize(width, height int) {
	//v = NewView(int64(width), int64(height))
	if len(*v) < height {
		for i := len(*v); i < height; i++ {
			*v = append((*v), make(CellList, width))
			for j := range (*v)[i] {
				(*v)[i][j] = newCell()
			}
		}
	} else {
		*v = (*v)[:height]
	}

	if len((*v)[0]) < width {
		for i := range *v {
			for j := len((*v)[i]); j < width; j++ {
				(*v)[i] = append((*v)[i], newCell())
			}
		}
	} else {
		for i := range *v {
			(*v)[i] = (*v)[i][:width]
		}
	}
}

func (v *View) GetCell(x, y int) *Cell {
	return (*v)[y][x]
}
