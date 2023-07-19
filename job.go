package myUtils

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type MyCron interface {
	AddFunc(spec string, cmd func()) (cron.EntryID, error)
	Remove(entryID cron.EntryID)
	StartCron() bool
	StopCron() bool
}

type job struct {
	log *Logger

	name   string
	status bool
	cron   *cron.Cron
}

func NewJob(jobName string, log *Logger) *job {
	return &job{
		name:   jobName,
		status: false,
		cron:   cron.New(cron.WithChain()),
	}
}

func (c *job) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	return c.cron.AddFunc(spec, cmd)
}

func (c *job) Remove(entryID cron.EntryID) {
	c.cron.Remove(entryID)
}

func (c *job) StartCron() bool {
	if !c.status {
		c.cron.Start()
		c.status = !c.status
		c.log.Info("定时任务已开启", zap.String("job", c.name))
	}
	return true
}

func (c *job) StopCron() bool {
	if c.status {
		c.cron.Stop()
		c.status = !c.status
		c.log.Info("定时任务已停止", zap.String("job", c.name))
	}
	return true
}
