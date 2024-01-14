package tty_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/saman3d/samtui/core/engine/tty"
)

// // devTty is an implementation of the Tty API based upon /dev/tty.
// type devTty struct {
// 	fd    int
// 	f     *os.File
// 	of    *os.File // the first open of /dev/tty
// 	saved *term.State
// 	sig   chan os.Signal
// 	cb    func()
// 	stopQ chan struct{}
// 	dev   string
// 	wg    sync.WaitGroup
// 	l     sync.Mutex
// }

// func (tty *devTty) Read(b []byte) (int, error) {
// 	return tty.f.Read(b)
// }

// func (tty *devTty) Write(b []byte) (int, error) {
// 	return tty.f.Write(b)
// }

// func (tty *devTty) Close() error {
// 	return tty.f.Close()
// }

// func (tty *devTty) Start() error {
// 	tty.l.Lock()
// 	defer tty.l.Unlock()

// 	// We open another copy of /dev/tty.  This is a workaround for unusual behavior
// 	// observed in macOS, apparently caused when a subshell (for example) closes our
// 	// own tty device (when it exits for example).  Getting a fresh new one seems to
// 	// resolve the problem.  (We believe this is a bug in the macOS tty driver that
// 	// fails to account for dup() references to the same file before applying close()
// 	// related behaviors to the tty.)  We're also holding the original copy we opened
// 	// since closing that might have deleterious effects as well.  The upshot is that
// 	// we will have up to two separate file handles open on /dev/tty.  (Note that when
// 	// using stdin/stdout instead of /dev/tty this problem is not observed.)
// 	var err error
// 	if tty.f, err = os.OpenFile(tty.dev, os.O_RDWR, 0); err != nil {
// 		return err
// 	}

// 	if !term.IsTerminal(tty.fd) {
// 		return errors.New("device is not a terminal")
// 	}

// 	_ = tty.f.SetReadDeadline(time.Time{})
// 	saved, err := term.MakeRaw(tty.fd) // also sets vMin and vTime
// 	if err != nil {
// 		return err
// 	}
// 	tty.saved = saved

// 	tty.stopQ = make(chan struct{})
// 	tty.wg.Add(1)
// 	go func(stopQ chan struct{}) {
// 		defer tty.wg.Done()
// 		for {
// 			select {
// 			case <-tty.sig:
// 				tty.l.Lock()
// 				cb := tty.cb
// 				tty.l.Unlock()
// 				if cb != nil {
// 					cb()
// 				}
// 			case <-stopQ:
// 				return
// 			}
// 		}
// 	}(tty.stopQ)

// 	signal.Notify(tty.sig, syscall.SIGWINCH)
// 	return nil
// }

// func (tty *devTty) Stop() error {
// 	tty.l.Lock()
// 	if err := term.Restore(tty.fd, tty.saved); err != nil {
// 		tty.l.Unlock()
// 		return err
// 	}
// 	_ = tty.f.SetReadDeadline(time.Now())

// 	signal.Stop(tty.sig)
// 	close(tty.stopQ)
// 	tty.l.Unlock()

// 	tty.wg.Wait()

// 	// close our tty device -- we'll get another one if we Start again later.
// 	_ = tty.f.Close()

// 	return nil
// }

// func (tty *devTty) WindowSize() (int, int, error) {
// 	w, h, err := term.GetSize(tty.fd)
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	if w == 0 {
// 		w, _ = strconv.Atoi(os.Getenv("COLUMNS"))
// 	}
// 	if w == 0 {
// 		w = 80 // default
// 	}
// 	if h == 0 {
// 		h, _ = strconv.Atoi(os.Getenv("LINES"))
// 	}
// 	if h == 0 {
// 		h = 25 // default
// 	}
// 	return w, h, nil
// }

// func (tty *devTty) NotifyResize(cb func()) {
// 	tty.l.Lock()
// 	tty.cb = cb
// 	tty.l.Unlock()
// }

// // NewDevTty opens a /dev/tty based Tty.
// func NewDevTty() (*devTty, error) {
// 	return NewDevTtyFromDev("/dev/tty")
// }

// // NewDevTtyFromDev opens a tty device given a path.  This can be useful to bind to other nodes.
// func NewDevTtyFromDev(dev string) (*devTty, error) {
// 	tty := &devTty{
// 		dev: dev,
// 		sig: make(chan os.Signal),
// 	}
// 	var err error
// 	if tty.of, err = os.OpenFile(dev, os.O_RDWR, 0); err != nil {
// 		return nil, err
// 	}
// 	tty.fd = int(tty.of.Fd())
// 	if !term.IsTerminal(tty.fd) {
// 		_ = tty.f.Close()
// 		return nil, errors.New("not a terminal")
// 	}
// 	if tty.saved, err = term.GetState(tty.fd); err != nil {
// 		_ = tty.f.Close()
// 		return nil, fmt.Errorf("failed to get state: %w", err)
// 	}
// 	return tty, nil
// }

func TestDevTty(t *testing.T) {
	ch := make(chan tty.Event, 1)

	tt, err := tty.NewTTY()
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println("tty", tt)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = tt.Start(ctx, ch)
	if err != nil {
		t.Fatal(err)
		return
	}

	go func() {
		for e := range ch {
			fmt.Println("event", e, e.(tty.KeyboardEvent).Key.Is(tty.SpecialKey_F12))
			fmt.Println("\r")
		}
	}()

	tt.Wait()
}

func TestSetCursor(t *testing.T) {
	tt, err := tty.NewTTY()
	if err != nil {
		t.Fatal(err)
		return
	}
	// clear screen
	fmt.Print("\033[2J")

	// move cursor on the screen
	w, h, err := tt.WindowSize()
	if err != nil {
		t.Fatal(err)
		return
	}
	for i := 1; i <= w; i++ {
		for j := 1; j <= h; j++ {
			tt.WritePos(i, j, []byte{fmt.Sprintf("%d", i)[0]})
			time.Sleep(1 * time.Millisecond)
		}
	}

}
