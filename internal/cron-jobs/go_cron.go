package cron_jobs

import (
	"time"

	"github.com/go-co-op/gocron"
)

// compile-time interface check
var _ CronJob = (*Cron)(nil)

type Cron struct {
	s *gocron.Scheduler
}

type CronOptions struct {
}

func NewCron() *Cron {
	return &Cron{
		s: gocron.NewScheduler(time.Local),
	}
}

func (c *Cron) AddJobWithInterval(interval any, cmd func()) error {
	_, err := c.s.Every(interval).Do(cmd)

	return err
}

func (c *Cron) AddJob(cmd func()) error {
	_, err := c.s.Do(cmd)

	return err
}

func (c *Cron) SetupJob(sch *gocron.Scheduler) *Cron {
	c.s = sch
	return c
}

func (c *Cron) StartWithBlocking() {
	c.s.StartBlocking()
}

func (c *Cron) StartAsync() {
	c.s.StartAsync()
}

func (c *Cron) IsRunning() bool {
	return c.s.IsRunning()
}

func (c *Cron) Stop() {
	c.s.Stop()
}
