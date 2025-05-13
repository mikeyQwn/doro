package bin

import (
	"io"
	"os"
	"sync"

	"github.com/faiface/beep"
	"github.com/mikeyQwn/doro/lib/input"
	"github.com/mikeyQwn/doro/lib/ui"
)

type AppState struct {
	wr            io.Writer
	ks            input.KeyStream
	streamer      beep.StreamSeekCloser
	isPlaying     bool
	isPlayingLock sync.Mutex
	cfg           *Config
}

func NewAppState(ks input.KeyStream, streamer beep.StreamSeekCloser, cfg *Config) *AppState {
	return &AppState{
		wr:       os.Stdout,
		ks:       ks,
		streamer: streamer,
		cfg:      cfg,
	}
}

func (s *AppState) NewWidget(update ui.UpdateFn) *ui.Widget {
	return ui.NewWidget(update).
		WithWriter(s.wr).
		EnableKeyHandling(s.ks)
}
