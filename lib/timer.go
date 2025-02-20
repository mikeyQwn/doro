package lib

import "time"

type Timer struct {
	start    time.Time
	duration time.Duration

	offs     time.Duration
	pausedAt time.Time
}

// Create a new `Timer` and starts it
func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		start:    time.Now(),
		duration: duration,

		offs:     0,
		pausedAt: time.Time{},
	}
}

func NewPaused(duration time.Duration) *Timer {
	t := NewTimer(duration)
	t.pausedAt = t.start

	return t
}

func (t *Timer) Pause() {
	if t.IsPaused() {
		return
	}

	t.pausedAt = time.Now()
}

func (t *Timer) Toggle() {
	if t.IsPaused() {
		t.Unpause()
	} else {
		t.Pause()
	}
}

func (t *Timer) Unpause() {
	if !t.IsPaused() {
		return
	}

	t.offs += time.Since(t.pausedAt)
	t.pausedAt = time.Time{}
}

func (t *Timer) IsPaused() bool {
	return (t.pausedAt != time.Time{})
}

func (t *Timer) Progress() float64 {
	elapsed := t.Elapsed()
	progress := float64(elapsed) / float64(t.duration)
	if progress > 1.0 {
		progress = 1.0
	}

	return progress
}

func (t *Timer) IsFinished() bool {
	return t.Elapsed() >= t.duration
}

// Returns time elapsed from the start, not counting pauses
func (t *Timer) Elapsed() time.Duration {
	now := time.Now()
	elapsed := now.Sub(t.start) - t.offs
	if (t.pausedAt != time.Time{}) {
		elapsed -= now.Sub(t.pausedAt)
	}

	return elapsed
}
