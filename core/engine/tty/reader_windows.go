//go:build windows

package tty

import (
	"fmt"
)

// --------------------
//    Windows Reader
// --------------------

type windowsReader struct {
	f *os.File
}

func newInputReader() (InputReader, error) {
	var tty windowsReader
	var err error

	// 1. Get the file descriptor of the TTY.
	tty.f, err = os.OpenFile("CONIN$", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	return &tty, nil
}

func (tty *windowsReader) Read(p []byte) (int, error) {
	return tty.f.Read(p)
}

func (tty *windowsReader) Close() error {
	return tty.f.Close()
}
