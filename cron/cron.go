package cron

import (
	cron "github.com/robfig/cron/v3"
	"time"
)

type Cron struct {
	*cron.Cron
}

// NewCron cron.New(cron.WithSeconds())
func NewCron(opts ...cron.Option) *Cron {
	return &Cron{Cron: cron.New(opts...)}
}

// AddDelayJob 添加定时任务
func (c *Cron) AddDelayJob(delay time.Duration, cmd func()) (cron.EntryID, error) {
	job := newWrapperCmd(cmd, c)
	id := c.Schedule(cron.Every(delay), job)
	f, ok := job.(interface{ SetId(id cron.EntryID) })
	if ok {
		f.SetId(id)
	}
	return id, nil
}

func newWrapperCmd(f func(), c *Cron) cron.Job {
	return &wrapperCmd{f: f, Cron: c}
}

type wrapperCmd struct {
	id cron.EntryID
	f  func()
	*Cron
}

func (w *wrapperCmd) Run() {
	w.f()
	w.Remove(w.id)
}

func (w *wrapperCmd) SetId(id cron.EntryID) {
	w.id = id
}
