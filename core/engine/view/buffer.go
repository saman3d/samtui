package view

import (
	"fmt"

	"github.com/saman3d/samtui/core/dom"
)

type Position struct {
	X int64
	Y int64
}

type Cell struct {
	Style   Style
	Content rune
	Element *dom.Element
	ZIndex  uint8
}

func newCell() *Cell {
	return &Cell{
		Content: ' ',
	}
}

func (c *Cell) Flush() {
	c.Style = Style{}
	c.Content = ' '
	c.Element = nil
}

func (c *Cell) String() string {
	return fmt.Sprintf("%s%c", c.Style, c.Content)
}

func (c *Cell) Bytes() []byte {
	return []byte(c.String())
}

type Style struct {
	Foreground int
	Background int
}

func NewStyle(fg, bg int) Style {
	return Style{
		Foreground: fg,
		Background: bg,
	}
}

func (s Style) String() string {
	if s.Background == 0 && s.Foreground == 0 {
		return fmt.Sprintf("\033[m")
	} else if s.Background == 0 {
		return fmt.Sprintf("\033[m\033[38;5;%vm", s.Foreground)
	} else if s.Foreground == 0 {
		return fmt.Sprintf("\033[m\033[48;5;%vm", s.Background)
	}
	return fmt.Sprintf("\033[38;5;%vm\033[48;5;%vm", s.Foreground, s.Background)
}

type AxisMask byte

const (
	AxisMask_X = 1 << iota
	AxisMask_Y
)
