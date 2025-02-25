package bin

import (
	"time"

	"github.com/mikeyQwn/doro/lib"
)

type Task struct {
	label string
	timer *lib.Timer
}

type Pomodoro struct {
	workTask  Task
	breakTask Task
	activeIdx int
}

func NewPomodoro(cfg *Config, isLong bool) *Pomodoro {
	pm := &Pomodoro{
		workTask: Task{
			label: focusedWorkLabel,
			timer: lib.NewTimer(cfg.focusedWorkDuration),
		},
		breakTask: Task{
			label: shortBreakLabel,
			timer: lib.NewPaused(cfg.shortBreakDuration),
		},
		activeIdx: 0,
	}
	if isLong {
		pm.breakTask.label = longBreakLabel
		pm.breakTask.timer = lib.NewPaused(cfg.longBreakDuration)
	}

	return pm
}

func (p *Pomodoro) WorkLabel() string {
	return p.workTask.label
}

func (p *Pomodoro) WorkProgress() float64 {
	return p.workTask.timer.Progress()
}

func (p *Pomodoro) WorkRunning() bool {
	return !p.workTask.timer.IsPaused() && p.workTask.timer.Elapsed() != 0
}

func (p *Pomodoro) WorkElapsed() time.Duration {
	return p.workTask.timer.Elapsed()
}

func (p *Pomodoro) WorkDuration() time.Duration {
	return p.workTask.timer.Duratoin()
}

func (p *Pomodoro) BreakLabel() string {
	return p.breakTask.label
}

func (p *Pomodoro) BreakProgress() float64 {
	return p.breakTask.timer.Progress()
}

func (p *Pomodoro) BreakRunning() bool {
	return !p.breakTask.timer.IsPaused() && p.breakTask.timer.Elapsed() != 0
}

func (p *Pomodoro) BreakElapsed() time.Duration {
	return p.breakTask.timer.Elapsed()
}

func (p *Pomodoro) BreakDuration() time.Duration {
	return p.breakTask.timer.Duratoin()
}

func (p *Pomodoro) Active() *Task {
	switch p.activeIdx {
	case 0:
		return &p.workTask
	case 1:
		return &p.breakTask
	}

	return nil
}

func (p *Pomodoro) Update() {
	task := p.Active()
	if task == nil {
		return
	}

	if !task.timer.IsFinished() {
		return
	}
	p.activeIdx += 1
	if task := p.Active(); task != nil {
		task.timer.Unpause()
	}
}

func (p *Pomodoro) IsPaused() bool {
	task := p.Active()
	if task == nil {
		return false
	}

	return task.timer.IsPaused()
}

func (p *Pomodoro) TogglePause() {
	task := p.Active()
	if task == nil {
		return
	}

	task.timer.Toggle()
}

func (p *Pomodoro) IsFinished() bool {
	return p.Active() == nil
}
