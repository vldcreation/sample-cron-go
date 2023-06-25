package cron_jobs

import "github.com/go-co-op/gocron"

type CronJob interface {
	SetupJob(sch *gocron.Scheduler) *Cron
	AddJob(cmd func()) error
}
