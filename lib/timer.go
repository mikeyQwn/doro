package lib

import "time"

type Timer struct {
	start    time.Time
	duration time.Duration

	offs     time.Duration
	pausedAt time.Time
}

// Creates a new `Timer` with a given duration and starts it
func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		start:    time.Now(),
		duration: duration,

		offs:     0,
		pausedAt: time.Time{},
	}
}

// Creates a new `Timer` with a given duration, but does not start it
func NewPaused(duration time.Duration) *Timer {
	t := NewTimer(duration)
	t.pausedAt = t.start

	return t
}

// Pauses the timer. Time spent between the `Pause()`
// and `Unpause()` calls does not go towards reaching
// timer's duration
func (t *Timer) Pause() {
	if t.IsPaused() {
		return
	}

	t.pausedAt = time.Now()
}

// Unpauses the timer
func (t *Timer) Unpause() {
	if !t.IsPaused() {
		return
	}

	t.offs += time.Since(t.pausedAt)
	t.pausedAt = time.Time{}
}

// Pauses a timer if it's unpaused or unpauses it
// if it's paused
func (t *Timer) Toggle() {
	if t.IsPaused() {
		t.Unpause()
	} else {
		t.Pause()
	}
}

// Returns `true` if the timer is paused
func (t *Timer) IsPaused() bool {
	return (t.pausedAt != time.Time{})
}

// Returns a float value in range [0.0, 1.0] that shows
// what percentage of timer's duration passed
func (t *Timer) Progress() float64 {
	elapsed := t.Elapsed()
	progress := float64(elapsed) / float64(t.duration)
	progress = min(1.0, progress)

	return progress
}

// Returns wether or not the timer is finished (e.g. reached it's duration)
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

	return max(elapsed, t.duration)
}

// Returns total timer duration
func (t *Timer) Duratoin() time.Duration {
	return t.duration
}
