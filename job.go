package main

import (
	"fmt"
	"time"

	"github.com/ricdeau/drop-test-db-job/internal"
	"github.com/ricdeau/drop-test-db-job/sql"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Job struct {
	id       string
	ctx      *internal.Context
	pipe     internal.Pipe
	schedule cron.Schedule
}

func (j *Job) Run() {
	err := j.pipe.
		Connect().
		Filter().
		Drop().
		Close()
	if err != nil {
		j.ctx.Logger.Error("Job finished with error.", zap.Error(err), zap.String("job-id", j.id))
	}
}

func NewPostgresJob(ctx *internal.Context, schedule cron.Schedule) *Job {
	id := jobId("postgres")
	return &Job{
		id:       id,
		ctx:      ctx,
		pipe:     sql.NewPostgresPipe(ctx, id),
		schedule: schedule,
	}
}

func jobId(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, time.Now().Format(time.RFC3339))
}
