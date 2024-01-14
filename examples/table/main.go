package main

import (
	"context"
	"fmt"
	"os"

	"github.com/saman3d/samtui/core/dom"
	"github.com/saman3d/samtui/core/engine"
	"github.com/saman3d/samtui/core/engine/tty"
)

func main() {
	app := NewApplication()

	app.Start()
}

type Application struct {
	eng      *engine.Engine
	num_rows int
	selected int
	modal    bool
}

func NewApplication() *Application {
	var a Application

	tf, err := os.Open("table.html")
	if err != nil {
		panic(err)
	}

	a.eng, err = engine.NewEngineFromTemplate(tf)
	if err != nil {
		panic(err)
	}

	return &a
}

func (a *Application) Start() {
	go a.eventWatcher()

	a.eng.Start(context.Background())
}

func (a *Application) eventWatcher() {
	var event tty.Event
	for event = range a.eng.PollEvent() {
		switch event.Type() {
		case tty.EventType_Keyboard:
			a.parseKeyboardEvent(event.(tty.KeyboardEvent))
		case tty.EventType_Resize:
			a.eng.Reload()
		}
	}
}

func (a *Application) parseKeyboardEvent(event tty.KeyboardEvent) {
	switch {
	case event.Is(tty.NewByteKey('c'), tty.Modifier_Ctrl):
		a.eng.Exit()
	case event.Is(tty.NewByteKey('i')):
		a.insertRow()
	case event.Is(tty.SpecialKey_Down):
		a.SelectNext()
	case event.Is(tty.SpecialKey_Up):
		a.SelectPrevious()
	case event.Is(tty.NewByteKey('\r')):
		a.ToggleModal(fmt.Sprintf("you have selected row %d", a.selected))
	case event.Is(tty.NewByteKey('q')):
		a.eng.Exit()
	}
}

func (a *Application) ToggleModal(s string) {
	if a.modal {
		a.HideModal()
		return
	}

	a.DrawModal(s)
}

func (a *Application) HideModal() {
	a.eng.GetElementByID("modal")[0].Remove()
	a.modal = false
}

func (a *Application) DrawModal(text string) {
	bnd := a.eng.GetElementByID("body")[0].Element().Boundry
	e := dom.MustParseElementFromString(fmt.Sprintf(`<div display="absolute" id="modal" z-index="3" height="10" width="40" background-color="0" left="%d" top="%d" border="true">%s</div>`, bnd.Width()/2-20, bnd.Height()/2-5, text))
	a.eng.GetElementByID("body")[0].PrependChild(e)
	a.modal = true
}

func (a *Application) SelectPrevious() {
	if a.num_rows == 0 || a.selected == 0 {
		return
	}

	tb := a.eng.GetElementByID("table-body")[0].Element()
	if a.selected == tb.State.ScrollY {
		a.ScrollTop()
	}

	tb.Children[a.selected].Attrs.BackGroundColor = 0
	tb.Children[a.selected].InheritChildrensAttr()
	a.eng.Update(tb.Children[a.selected])
	a.selected--
	tb.Children[a.selected].Attrs.BackGroundColor = 56
	tb.Children[a.selected].InheritChildrensAttr()
	a.eng.Update(tb.Children[a.selected])
}

func (a *Application) SelectNext() {
	if a.num_rows == 0 || a.selected >= a.num_rows-1 {
		return
	}

	tb := a.eng.GetElementByID("table-body")[0].Element()
	if a.selected+2 > tb.State.ScrollY+tb.Boundry.Height() {
		a.ScrollDown()
	}

	tb.Children[a.selected].Attrs.BackGroundColor = 0
	tb.Children[a.selected].InheritChildrensAttr()
	a.eng.Update(tb.Children[a.selected])
	a.selected++
	tb.Children[a.selected].Attrs.BackGroundColor = 56
	tb.Children[a.selected].InheritChildrensAttr()
	a.eng.Update(tb.Children[a.selected])
}

func (a *Application) ScrollDown() {
	tb := a.eng.GetElementByID("table-body")[0].Element()
	if a.num_rows > tb.Boundry.Height() {
		if a.num_rows > tb.State.ScrollY+tb.Boundry.Height() {
			tb.State.ScrollBy(0, 1)
			a.eng.Update(tb)
		}
	}
}

func (a *Application) ScrollTop() {
	tb := a.eng.GetElementByID("table-body")[0].Element()
	if tb.State.ScrollY >= 1 {
		tb.State.ScrollBy(0, -1)
		a.eng.Update(tb)
	}
}

func (a *Application) insertRow() {
	tb := a.eng.GetElementByID("table-body")[0]
	d := len(tb.Element().Children)
	tb.AppendChild(dom.MustParseElementFromString(fmt.Sprintf(`
                <trow display="flex" height="1">
                    <p>row %d</p>
                    <p>row %d</p>
                    <p>row %d</p>
                    <p>row %d</p>
                </trow>
		`, d, d, d, d)))
	a.num_rows++
}
