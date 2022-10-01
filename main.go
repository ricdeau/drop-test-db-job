package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/ricdeau/drop-test-db-job/context"
	"github.com/ricdeau/drop-test-db-job/plugins"
	"github.com/robfig/cron/v3"
)

type Command struct {
	AgeThreshold  time.Duration
	DbType        string
	DbNamePattern string
	DSN           string
	Schedule      string
}

func main() {
	command := kong.Parse(&Command{}, kong.Bind(context.New()))
	command.FatalIfErrorf(command.Run())
}

func (c *Command) AfterApply(ctx *context.Context) (err error) {
	ctx.AgeThreshold = c.AgeThreshold
	ctx.DSN = c.DSN
	ctx.DbNamePattern = c.DbNamePattern

	ctx.Schedule, err = cron.ParseStandard(c.Schedule)
	if err != nil {
		return fmt.Errorf("parse schedule: %v", err)
	}
	ctx.Plugin, err = plugins.GetPlugin(c.DbType)
	if err != nil {
		return fmt.Errorf("get plugin for %s: %v", c.DbType, err)
	}

	return nil
}

func (c *Command) Run(ctx *Context) error {
	newCron := cron.New()
	for _, j := range ctx.Jobs {
		newCron.Schedule(j.GetSchedule(), j)
	}
	sig := make(chan os.Signal, 1)
	go func() {
		<-sig
		<-newCron.Stop().Done()
		ctx.Cancel()
	}()

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP)

	<-ctx.Done()

	return nil
}
