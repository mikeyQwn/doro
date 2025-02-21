package bin

import (
	"io"
	"os"

	"github.com/mikeyQwn/doro/lib/input"
	"github.com/mikeyQwn/doro/lib/ui"
)

type AppState struct {
	wr  io.Writer
	ks  input.KeyStream
	cfg *Config
}

func NewAppState(ks input.KeyStream, cfg *Config) *AppState {
	return &AppState{
		wr:  os.Stdout,
		ks:  ks,
		cfg: cfg,
	}
}

func (s *AppState) NewWidget(update ui.UpdateFn) *ui.Widget {
	return ui.NewWidget(update).
		WithWriter(s.wr).
		EnableKeyHandling(s.ks)
}
