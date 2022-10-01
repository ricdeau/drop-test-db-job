package main

import (
	"context"

	"github.com/ricdeau/drop-test-db-job/internal"
)

type Context struct {
	context.Context
	Cancel context.CancelFunc

	Jobs []internal.Job
}
