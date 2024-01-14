package tty

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/term"
)

type TTY struct {
	inpreader InputReader
	inpparser InputParser
	wg        *sync.WaitGroup
	cursor    [2]int
	sigwinch  chan os.Signal
}

func NewTTY() (*TTY, error) {
	var err error
	tty := &TTY{
		wg:       &sync.WaitGroup{},
		cursor:   [2]int{},
		sigwinch: make(chan os.Signal),
	}

	tty.inpreader, err = newInputReader()
	if err != nil {
		return nil, err
	}

	signal.Notify(tty.sigwinch, syscall.SIGWINCH)

	tty.inpparser = newInputParser()

	tty.DisableCursor()

	return tty, nil
}

func (t *TTY) Wait() {
	t.wg.Wait()
}

func (t *TTY) Watch(ctx context.Context, echan chan Event) error {
	inpchan := make(chan []byte)
	go t.readChan(inpchan)
	for {
		select {
		case <-ctx.Done():
			t.Close()
			return ctx.Err()
		case b := <-inpchan:
			// 1. Parse the input.
			input := t.inpparser.Parse(b)

			// 2. Send input events to the channel.
			for _, i := range input {
				echan <- i
			}

		case <-t.sigwinch:
			w, h, _ := t.WindowSize()
			echan <- ResizeEvent{
				Width:  w,
				Height: h,
			}
		}
	}
}

func (t *TTY) readChan(ch chan []byte) {
	for {
		// 1. Read input from the TTY.
		b := make([]byte, 10)
		_, err := t.inpreader.Read(b)
		if err != nil {
			fmt.Println(err)
			// return
		}

		// 2. Send input to the channel.
		ch <- b
	}

}

func (t *TTY) Start(ctx context.Context, ch chan Event) error {
	t.wg.Add(1)
	go t.Watch(ctx, ch)
	return nil
}

func (t *TTY) Close() error {
	err := t.inpreader.Close()
	if err != nil {
		return err
	}
	return nil
}

func (t *TTY) WindowSize() (int, int, error) {
	t.inpreader, _ = newInputReader()
	return term.GetSize(int(t.inpreader.Fd()))
}

func (t *TTY) Clear() {
	t.inpreader.Write([]byte("\033[2J"))
}

func (t *TTY) WritePos(x, y int, b []byte) (int, error) {
	if t.cursor[0] != x || t.cursor[1] != y {
		t.SetPos(x, y)
	}
	return t.inpreader.Write(b)
}

func (t *TTY) Write(b []byte) (int, error) {
	b = append(b, []byte(ResetEscape)...)
	t.setCursor(t.cursor[0]+len(b), t.cursor[1])
	return t.inpreader.Write(b)
}

func (t *TTY) DisableCursor() {
	t.inpreader.Write([]byte("\033[?25l"))
}

func (t *TTY) setCursor(x, y int) {
	t.cursor[0] = x
	t.cursor[1] = y
}

func (t *TTY) SetPos(x, y int) {
	t.cursor[0], t.cursor[1] = x, y
	t.inpreader.Write([]byte("\033[" + fmt.Sprint(y+1) + ";" + fmt.Sprint(x+1) + "H"))
}
