package main

import (
	"github.com/ricdeau/drop-test-db-job/internal"
	"github.com/ricdeau/jsdur"
)

type Config struct {
	Jobs []struct {
		Type             internal.DatabaseType  `json:"type"`
		ConnectionString string                 `json:"connectionString"`
		TimestampType    internal.TimestampType `json:"timestampType"`
		MaxAge           jsdur.Duration         `json:"maxAge"`
		Timeout          jsdur.Duration         `json:"timeout"`
		Schedule         internal.Schedule      `json:"schedule"`
	} `json:"jobs"`
}
