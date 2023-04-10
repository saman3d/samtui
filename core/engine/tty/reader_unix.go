//go:build unix

package tty

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

// --------------------
//     Unix Reader
// --------------------

type unixReader struct {
	fd    uintptr
	f     *os.File
	saved *term.State
}

func newInputReader() (InputReader, error) {
	var tty unixReader
	var err error

	// 1. Get the file descriptor of the TTY.
	tty.f, err = os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	tty.fd = tty.f.Fd()
	_, err = fcntl(int(tty.fd), syscall.F_SETFL, syscall.O_ASYNC|syscall.O_NONBLOCK)

	if !term.IsTerminal(int(tty.f.Fd())) {
		return nil, fmt.Errorf("not a terminal")
	}

	// 2. Save the current state of the TTY
	// and set the terminal to raw mode.
	tty.saved, err = term.MakeRaw(int(tty.fd))
	if err != nil {
		return nil, err
	}
	tty.saved, err = term.GetState(int(tty.fd))
	if err != nil {
		return nil, err
	}

	return &tty, nil
}

func (tty *unixReader) Read(p []byte) (int, error) {
	return tty.f.Read(p)
}

func (tty *unixReader) Write(p []byte) (int, error) {
	return tty.f.Write(p)
}

func (tty *unixReader) Close() error {
	term.Restore(int(tty.fd), tty.saved)
	return tty.f.Close()
}

func (tty *unixReader) Fd() uintptr {
	return tty.fd
}

func fcntl(fd int, cmd int, arg int) (val int, err error) {
	r, _, e := syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), uintptr(cmd),
		uintptr(arg))
	val = int(r)
	if e != 0 {
		err = e
	}
	return
}
