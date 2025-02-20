package bin

import (
	"io"
	"os"

	"github.com/mikeyQwn/doro/lib"
	"github.com/mikeyQwn/doro/lib/input"
	"github.com/mikeyQwn/doro/lib/ui"
)

type AppState struct {
	w   *lib.BatchWriter
	wr  io.Writer
	ks  input.KeyStream
	cfg Config
}

func NewAppState(ks input.KeyStream) *AppState {
	return &AppState{
		w:   lib.NewBatchWriter(os.Stdout),
		wr:  os.Stdout,
		ks:  ks,
		cfg: Config{},
	}
}

func (s *AppState) NewWidget(update ui.UpdateFn) *ui.Widget {
	return ui.NewWidget(update).
		WithWriter(s.wr).
		EnableKeyHandling(s.ks)
}
