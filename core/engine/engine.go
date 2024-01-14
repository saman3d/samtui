package engine

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/bep/debounce"
	"github.com/saman3d/samtui/core/dom"
	"github.com/saman3d/samtui/core/engine/tty"
	"github.com/saman3d/samtui/core/engine/view"
)

var ids map[string][]*dom.Element

type Engine struct {
	DOM     *dom.Document
	View    View
	TTY     TTY
	Layouts map[LayoutType]Layout

	eventch     chan tty.Event
	renderstack RenderStack

	cancel func()
	dbnc   func(func())
}

func NewEngineFromTemplate(tf io.Reader) (*Engine, error) {
	dm, err := dom.NewDocumentFromReader(tf)
	if err != nil {
		return nil, err
	}

	t, err := tty.NewTTY()
	if err != nil {
		return nil, err
	}

	width, height, err := t.WindowSize()
	if err != nil {
		return nil, err
	}

	v := view.NewView(int64(width), int64(height))

	dm.Body.Boundry = v.Boundry()

	renderstack := newRenderStack()

	e := &Engine{
		DOM:  dm,
		TTY:  t,
		View: v,
		Layouts: map[LayoutType]Layout{
			LayoutType_Flex:     newFlexLayout(v, renderstack),
			LayoutType_Block:    newBlockLayout(v, renderstack),
			LayoutType_Absolute: newAbsoluteLayout(v, renderstack),
		},
		renderstack: renderstack,
		eventch:     make(chan tty.Event, 10),
		dbnc:        debounce.New(time.Millisecond * 100),
	}

	e.generateIDsMap()

	return e, nil
}

func (e *Engine) reload() {
	e.TTY.Clear()
	w, h, _ := e.TTY.WindowSize()
	e.View.Resize(w, h)
	e.DOM.Body.Boundry = e.View.Boundry()
	e.renderstack.Push(e.DOM.Body)
}

func (e *Engine) Reload() {
	e.dbnc(e.reload)
}

func (e *Engine) generateIDsMap() {
	ids = make(map[string][]*dom.Element)
	e.generateIDsMapRecursive(e.DOM.Body)
}

func (e *Engine) generateIDsMapRecursive(el *dom.Element) {
	if el.Attrs.ID != "" {
		ids[el.Attrs.ID] = append(ids[el.Attrs.ID], el)

	}

	for _, child := range el.Children {
		e.generateIDsMapRecursive(child)
	}
}

func (e *Engine) GetElementByID(id string) []Elementor {
	elem, ok := ids[id]
	if !ok {
		return nil
	}

	elmtrs := make([]Elementor, len(elem))
	for i := range elem {
		elmtrs[i] = newElementUpdater(elem[i], e)
	}

	return elmtrs
}

func (e *Engine) Start(cp context.Context) error {
	var ctx context.Context
	ctx, e.cancel = context.WithCancel(cp)

	e.TTY.Clear()

	go e.TTY.Watch(ctx, e.eventch)

	go e.Render(ctx)

	e.renderstack.Push(e.DOM.Body)

	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (e *Engine) Render(ctx context.Context) {
	for range time.Tick(time.Millisecond * 1) {
		if e.renderstack.Len() == 0 {
			continue
		}
		var bnd *dom.Boundry
		for e.renderstack.Len() != 0 {
			bl := e.renderstack.Pop()
			if bl == nil {
				continue
			}
			e.renderElement(ctx, bl)
			if bnd == nil {
				bnd = &dom.Boundry{}
				*bnd = bl.Boundry
			} else if !bnd.Circumscribes(bl.Boundry) {
				*bnd = bnd.Sum(bl.Boundry)
			}
		}
		e.renderBoundry(*bnd)
	}
}

func (e *Engine) renderBoundry(bndr dom.Boundry) {
	for i := bndr.FirstY; i < bndr.SecondY; i++ {
		e.TTY.SetPos(bndr.FirstX, i)
		for j := bndr.FirstX; j < bndr.SecondX; j++ {
			e.TTY.Write(e.View.GetCell(j, i).Bytes())
		}
	}
}

func (e *Engine) renderElement(ctx context.Context, el *dom.Element) error {
	var err error
	switch el.Attrs.Display {
	case dom.Display_Flex:
		err = e.Layouts[LayoutType_Flex].Layout(ctx, el, el.Boundry)
		if err != nil {
			return err
		}
	case dom.Display_Block:
		err = e.Layouts[LayoutType_Block].Layout(ctx, el, el.Boundry)
		if err != nil {
			return err
		}
	case dom.Display_Absolute:
		err = e.Layouts[LayoutType_Absolute].Layout(ctx, el, el.Boundry)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) Update(el *dom.Element) {
	e.renderstack.Push(el)
}

func (e *Engine) PollEvent() <-chan tty.Event {
	return e.eventch
}

func (e *Engine) Exit() {
	e.TTY.Close()
	e.cancel()
}

type TTY interface {
	SetPos(x, y int)
	WindowSize() (int, int, error)
	Watch(ctx context.Context, ch chan tty.Event) error
	WritePos(x, y int, b []byte) (int, error)
	Write(b []byte) (int, error)
	Clear()
	Close() error
}

type View interface {
	Width() int64
	Height() int64
	Resize(width, height int)
	Boundry() dom.Boundry
	Flush()
	ClearBoundry(bndr dom.Boundry)
	PrintString(x, y, fg, bg int, zindx uint8, s string)
	PrintRune(x, y, fg, bg int, zindx uint8, r rune)
	PrintRuneRepeat(x, y, fg, bg, n int, zindx uint8, axis view.AxisMask, r rune)
	Slice(x, y, l int) view.CellList
	GetCell(x, y int) *view.Cell
}

type Elementor interface {
	Element() *dom.Element
	AppendChild(*dom.Element)
	PrependChild(*dom.Element)
	Remove()
	Update()
}

type ElementUpdater struct {
	el  *dom.Element
	eng *Engine
}

func newElementUpdater(el *dom.Element, eng *Engine) *ElementUpdater {
	return &ElementUpdater{
		el:  el,
		eng: eng,
	}
}

func (eu *ElementUpdater) Element() *dom.Element {
	return eu.el
}

func (eu *ElementUpdater) Update() {
	eu.eng.renderstack.Push(eu.el)
}

func (eu *ElementUpdater) AppendChild(el *dom.Element) {
	el.Parent = eu.el
	if el.Attrs.ID != "" {
		if l, ok := ids[el.Attrs.ID]; ok {
			l = append(l, el)
		} else {
			ids[el.Attrs.ID] = make([]*dom.Element, 0)
			ids[el.Attrs.ID] = append(ids[el.Attrs.ID], el)
		}
	}
	eu.el.AppendChild(el)
	eu.eng.renderstack.Push(eu.el)
}

func (eu *ElementUpdater) PrependChild(el *dom.Element) {
	el.Parent = eu.el
	if el.Attrs.ID != "" {
		if l, ok := ids[el.Attrs.ID]; ok {
			l = append(l, el)
		} else {
			ids[el.Attrs.ID] = make([]*dom.Element, 0)
			ids[el.Attrs.ID] = append(ids[el.Attrs.ID], el)
		}
	}
	eu.el.Children = append([]*dom.Element{el}, eu.el.Children...)
	eu.eng.renderstack.Push(eu.el)
}

func (eu *ElementUpdater) Remove() {
	eu.eng.View.ClearBoundry(eu.el.Boundry)
	if eu.el.Attrs.ID != "" {
		delete(ids, eu.el.Attrs.ID)
	}
	for i, ch := range eu.el.Parent.Children {
		if ch == eu.el {
			eu.el.Parent.Children = append(eu.el.Parent.Children[:i], eu.el.Parent.Children[i+1:]...)
			eu.eng.renderstack.Push(eu.el.Parent)
			return
		}
	}
}

type RenderStack interface {
	Push(el *dom.Element)
	Pop() *dom.Element
	Len() int
}

type renderStack struct {
	mu       sync.Mutex
	elements []*dom.Element
}

func newRenderStack() RenderStack {
	return &renderStack{}
}

func (rs *renderStack) Push(el *dom.Element) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	// remove last if element is already in stack
	for i, e := range rs.elements {
		if e == el {
			rs.elements = append(rs.elements[:i], rs.elements[i+1:]...)
			break
		}
	}
	rs.elements = append(rs.elements, el)
	copy(rs.elements[1:], rs.elements)
	rs.elements[0] = el
}

func (rs *renderStack) Pop() *dom.Element {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if len(rs.elements) == 0 {
		return nil
	}
	el := rs.elements[len(rs.elements)-1]
	rs.elements = rs.elements[:len(rs.elements)-1]
	return el
}

func (rs *renderStack) Len() int {
	return len(rs.elements)
}
