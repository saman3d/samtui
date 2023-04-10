package engine_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/saman3d/samtui/core/dom"
	"github.com/saman3d/samtui/core/engine"
	"github.com/saman3d/samtui/core/engine/tty"
)

var SimpleHtmlTemplate = `<!DOCTYPE html>
<html>
  <head>
	<title>Simple HTML Template</title>
  </head>
  <body display="flex" flex-direction="column">
    <div height="5" border="true"></div>
    <div height="5" border="true"></div>
    <div height="5" border="true"></div>
    <div height="5" border="true"></div>
    <div height="5" border="true"></div>
    <div height="5" border="true"></div>
    <div height="5" border="true"></div>
    <div height="5" border="true"></div>
    <div height="5" border="true"></div>
    <div height="5" border="true"></div>
    <div height="5" border="true"></div>
    <div height="5" border="true"></div>
  </body>
</html>`

var eng *engine.Engine

func TestEngineSimple(t *testing.T) {
	var err error
	f, err := os.Open("template_test.html")
	if err != nil {
		t.Fatal(err)
	}

	eng, err = engine.NewEngineFromTemplate(f)
	if err != nil {
		t.Fatal(err)
	}

	go mm()

	eng.Start(context.Background())
}

func mm() {
	element := eng.GetElementByID("table_body")

	for i := 0; i < 100; i++ {
		elem, _ := dom.NewElementFromString(fmt.Sprintf(`<div height="1" display="flex">
						<p id="pid" width="10">PID</p>
						<p id="cpu" width="10">CPU</p>
						<p id="mem" width="10">MEM</p>
						<p id="name">%d</p>
		</div>`, i))
		element[0].Element().Children = append(element[0].Element().Children, elem)
		// element[0].AppendChild(elem)
	}
	element[0].Update()
	var i int
	element[0].Element().Children[i].Attrs.BackGroundColor = 23
	element[0].Element().Children[i].Attrs.Color = 3
	element[0].Element().Children[i].InheritChildrensAttr()
	element[0].Update()
	for ev := range eng.PollEvent() {
		switch ev.Type() {
		case tty.EventType_Keyboard:
			ev := ev.(tty.KeyboardEvent)
			switch {
			case ev.Is(tty.NewByteKey('c'), tty.Modifier_Ctrl):
				eng.Exit()
			case ev.Is(tty.NewByteKey('q')):
				eng.Exit()
			case ev.Is(tty.SpecialKey_Down, tty.Modifier_None):
				i++
				element[0].Element().Children[i].Attrs.BackGroundColor = 23
				element[0].Element().Children[i].Attrs.Color = 3
				element[0].Element().Children[i].InheritChildrensAttr()
				element[0].Element().Children[i-1].Content = fmt.Sprintf("%d", 0)
				element[0].Element().Children[i-1].Attrs.BackGroundColor = 0
				element[0].Element().Children[i-1].Attrs.Color = 0
				element[0].Element().Children[i-1].InheritChildrensAttr()
				eng.Update(element[0].Element().Children[i])
				eng.Update(element[0].Element().Children[i-1])
			case ev.Is(tty.SpecialKey_Up, tty.Modifier_None):
				i--
				element[0].Element().Children[i].Content = fmt.Sprintf("%d", 23)
				element[0].Element().Children[i].Attrs.BackGroundColor = 23
				element[0].Element().Children[i].Attrs.Color = 3
				element[0].Element().Children[i].InheritChildrensAttr()
				element[0].Element().Children[i+1].Content = fmt.Sprintf("%d", 0)
				element[0].Element().Children[i+1].Attrs.BackGroundColor = 0
				element[0].Element().Children[i+1].Attrs.Color = 0
				element[0].Element().Children[i+1].InheritChildrensAttr()
				eng.Update(element[0].Element().Children[i])
				eng.Update(element[0].Element().Children[i+1])
			}

		}

	}
}
