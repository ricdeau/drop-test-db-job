package internal

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Context struct {
	context.Context
	ConnectionString string
	TimestampType    TimestampType
	MaxAge           time.Duration
	Logger           *zap.Logger
}
