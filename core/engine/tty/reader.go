package tty

import (
	"io"
)

// --------------------
//    Input Reader
// --------------------

type InputReader interface {
	io.ReadWriteCloser
	Fd() uintptr
}
