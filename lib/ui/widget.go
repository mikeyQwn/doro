package ui

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/mikeyQwn/doro/lib/terminal"
)

// Widget event handlers
type KeyHandler func(terminal.Key)

// `UpdateFn` is called when the widget is run and then after every event
// Returns widget "view" as a slice of lines
// Returning nil indicates that the widget should close
// Note that number of lines returned after every update should be the same
type UpdateFn func(f *Formatter) []string

// A terminal widget
// Note that widgets are not indended to run in parallel
type Widget struct {
	keyStream    terminal.KeyStream
	keyToHandler map[terminal.Key]KeyHandler
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
		w:            bufio.NewWriter(os.Stdout),
		n:            -1,
	}
}

// Attaches `keySteam` to the widget, enabling key event handling
func (w *Widget) EnableKeyHandling(keyStream terminal.KeyStream) *Widget {
	w.keyStream = keyStream
	return w
}

// Adds a handler that is triggered when a key is pressed
func (w *Widget) AddKeyHandler(handler KeyHandler, keys ...terminal.Key) *Widget {
	if w.keyToHandler == nil {
		w.keyToHandler = make(map[terminal.Key]KeyHandler, len(keys))
	}
	for _, k := range keys {
		w.keyToHandler[k] = handler
	}
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
		if w.n != 1 {
			w.w.WriteString(terminal.ESCAPE_CR_LF)
			w.w.Flush()
		}
	}()

	if shouldExit, err := w.triggerUpdate(); shouldExit || err != nil {
		return err
	}

	return w.mainloop()
}

func (w *Widget) mainloop() error {
	if w.keyStream != nil && w.keyToHandler != nil {
		for key := range w.keyStream {
			h, ok := w.keyToHandler[key]
			if !ok {
				continue
			}

			h(key)
			if shouldExit, err := w.triggerUpdate(); shouldExit || err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *Widget) triggerUpdate() (bool, error) {
	fmt, err := NewFormatter()
	if err != nil {
		return false, err
	}

	lines := w.update(fmt)
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

	if _, err = w.w.WriteString(strings.Join(lines,
		terminal.ESCAPE_CR_LF+terminal.ESCAPE_ERASE_LINE)); err != nil {
		return false, err
	}

	if _, err = w.w.WriteString(terminal.ESCAPE_CARRIAGE_RETURN); err != nil {
		return false, err
	}

	return false, w.w.Flush()
}
