package tty

import (
	"fmt"

	"github.com/saman3d/samtui/core/engine/view"
)

func PositionToEscapeCode(pos view.Position) string {
	return fmt.Sprintf("\033[%d;%dH", pos.Y, pos.X)
}
