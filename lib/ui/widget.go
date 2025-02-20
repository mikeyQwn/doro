package ui

import (
	"bufio"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mikeyQwn/doro/lib/ansi"
	"github.com/mikeyQwn/doro/lib/input"
	"github.com/mikeyQwn/doro/lib/terminal"
)

// Widget event handlers
type KeyHandler func(input.Key)
type TimedHandler func()

// `UpdateFn` is called when the widget is run and then after every event
// Returns widget "view" as a slice of lines and `shouldClose` boolean
// Truthy `shouldClose` indicates that the widget is supposed to be done after the next update
// Note that number of lines returned after every update should be the same
type UpdateFn func(*Formatter) ([]string, bool)

// A terminal widget
// Note that widgets are not indended to run in parallel
type Widget struct {
	keyStream    input.KeyStream
	keyToHandler map[input.Key]KeyHandler
	ticker       *time.Ticker
	timedHandler TimedHandler
	update       UpdateFn
	w            *bufio.Writer
	n            int
}

// Creates a new widget with an update function
func NewWidget(update UpdateFn) *Widget {
	return &Widget{
		keyStream:    nil,
		keyToHandler: nil,
		update:       update,
		ticker:       &time.Ticker{},
		timedHandler: nil,
		w:            bufio.NewWriter(os.Stdout),
		n:            -1,
	}
}

// Attaches `keySteam` to the widget, enabling key event handling
func (w *Widget) EnableKeyHandling(keyStream input.KeyStream) *Widget {
	w.keyStream = keyStream
	return w
}

// Adds a handler that is triggered when a key is pressed
func (w *Widget) AddKeyHandler(handler KeyHandler, keys ...input.Key) *Widget {
	if w.keyToHandler == nil {
		w.keyToHandler = make(map[input.Key]KeyHandler, len(keys))
	}
	for _, k := range keys {
		w.keyToHandler[k] = handler
	}
	return w
}

// Adds a handler that fires periodically every `duration`
func (w *Widget) AddTimedHandler(handler TimedHandler, duration time.Duration) *Widget {
	w.ticker = time.NewTicker(duration)
	w.timedHandler = handler
	return w
}

// Replaces the default writer (stdout) with `writer`
func (w *Widget) WithWriter(writer io.Writer) *Widget {
	w.w = bufio.NewWriter(writer)
	return w
}

// Runs the widget until an update returns nil
// Blocks until completion
func (w *Widget) Run() error {
	defer func() {
		_, _ = w.w.WriteString(ansi.CARRIAGE_RETURN + ansi.LINE_FEED)
		_ = w.w.Flush()
	}()

	if shouldExit, err := w.triggerUpdate(); shouldExit || err != nil {
		return err
	}

	return w.mainloop()
}

func (w *Widget) mainloop() error {
	if (w.keyStream == nil || w.keyToHandler == nil) &&
		(*w.ticker == time.Ticker{} || w.timedHandler == nil) {
		return nil
	}

	for {
		select {
		case key := <-w.keyStream:
			h, ok := w.keyToHandler[key]
			if !ok {
				continue
			}

			h(key)
			if shouldExit, err := w.triggerUpdate(); shouldExit || err != nil {
				return err
			}
		case <-w.ticker.C:
			w.timedHandler()
			if shouldExit, err := w.triggerUpdate(); shouldExit || err != nil {
				return err
			}
		}
	}
}

func (w *Widget) triggerUpdate() (bool, error) {
	fmt, err := NewFormatter()
	if err != nil {
		return false, err
	}

	lines, shouldClose := w.update(fmt)
	if lines == nil {
		return true, nil
	}

	if w.n == -1 {
		w.n = len(lines)

		if _, err := w.w.WriteString(
			terminal.DownLF(uint(w.n)) +
				terminal.Up(1),
		); err != nil {
			return false, err
		}
	}

	if w.n != 1 {
		if _, err = w.w.WriteString(terminal.Up(uint(w.n - 1))); err != nil {
			return false, err
		}
	}

	formattedWidget := strings.Join(lines,
		ansi.CARRIAGE_RETURN+ansi.LINE_FEED+ansi.ERASE_LINE,
	) + ansi.CARRIAGE_RETURN
	if _, err = w.w.WriteString(formattedWidget); err != nil {
		return false, err
	}

	return shouldClose, w.w.Flush()
}
