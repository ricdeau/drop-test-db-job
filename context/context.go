package context

import (
	"context"
	"time"

	"github.com/ricdeau/drop-test-db-job/plugins"
	"github.com/robfig/cron/v3"
)

type Context struct {
	context.Context
	Cancel context.CancelFunc

	AgeThreshold   time.Duration
	Plugin         *plugins.Plugin
	DbNamePattern  string
	TimestampPlace string
	DSN            string
	Schedule       cron.Schedule
}

func New() *Context {
	return &Context{
		Context: context.Background(),
	}
}

func (c *Context) Run() {
	err := c.Plugin.Connect(c.DSN)
	if err != nil {
		return
	}
	defer c.Plugin.Close()

	results, err := c.Plugin.GetOldDbs(c, c.DbNamePattern, c.TimestampPlace, c.AgeThreshold)
	if err != nil {
		return
	}

	for result := range results {
		if err = result.Err; err != nil {
			continue
		}

		if err = c.Plugin.DropDb(c, result.DbName); err != nil {
			continue
		}
	}
}
