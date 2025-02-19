package bin

import (
	"os"

	"github.com/mikeyQwn/doro/lib"
	"github.com/mikeyQwn/doro/lib/terminal"
)

type AppState struct {
	w   *lib.BatchWriter
	ks  terminal.KeyStream
	cfg Config
}

func NewAppState(ks terminal.KeyStream) *AppState {
	return &AppState{
		w:   lib.NewBatchWriter(os.Stdout),
		ks:  ks,
		cfg: Config{},
	}
}
