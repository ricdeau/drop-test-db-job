package internal

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
)

type Job interface {
	cron.Job
	GetSchedule() cron.Schedule
}

type Pipe interface {
	Connect() Connected
}

type Connected interface {
	Filter() Filtered
}

type Filtered interface {
	Drop() Closer
}

type Closer interface {
	Close() error
}

type DatabaseFilter interface {
	FilterDatabases(ctx context.Context, createdBefore time.Time) ([]string, error)
}

type DatabaseDropper interface {
	Drop(ctx context.Context, name string) error
}
