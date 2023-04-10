package engine

import (
	"fmt"
	"testing"
)

func TestRenderMan(t *testing.T) {
	rm := newRenderMan(10, RenderFlag_Flex)

	head := 0
	for f, d, ok := rm.Next(head, RenderFlag_Width|RenderFlag_MaxWidth); ok; f, d, ok = rm.Next(head, RenderFlag_Width|RenderFlag_MaxWidth) {
		fmt.Println(f, d)
		head++
	}
}
