package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ricdeau/drop-test-db-job/internal"
	"go.uber.org/zap"
)

const (
	PostgresPrefixPattern  = `^\d{13}[a-zA-Z0-9\-._]+$`
	PostgresPostfixPattern = `^[a-zA-Z0-9\-._]+\d{13}$`
)

type PostgresPipe struct {
	ctx           context.Context
	err           error
	dsn           string
	timestampType internal.TimestampType
	createdBefore time.Time
	getLogger     func() *zap.SugaredLogger
}

func NewPostgresPipe(ctx *internal.Context, id string) *PostgresPipe {
	return &PostgresPipe{
		ctx:           ctx.Context,
		dsn:           ctx.ConnectionString,
		timestampType: ctx.TimestampType,
		createdBefore: time.Now().Add(-ctx.MaxAge),
		getLogger: func() *zap.SugaredLogger {
			return ctx.Logger.
				Named("postgres_pipe").
				Sugar().
				With("job-id", id)
		},
	}
}

func (p *PostgresPipe) Connect() internal.Connected {
	db, err := sqlx.Connect(internal.DatabasePostgres.DriverName(), p.dsn)
	if err != nil {
		err = fmt.Errorf("connect: %v", err)
	}

	logger := p.getLogger()
	logger.Info("Connected")

	return &postgresConnected{
		PostgresPipe: p,
		db:           db,
		logger:       logger,
	}
}

type postgresConnected struct {
	*PostgresPipe
	logger *zap.SugaredLogger
	db     *sqlx.DB
}

func (p *postgresConnected) Filter() internal.Filtered {
	if p.err != nil {
		return &dropperCloser{err: p.err}
	}

	databases, err := postgresFilter(p.db, p.timestampType).FilterDatabases(p.ctx, p.createdBefore)
	if err != nil {
		p.err = fmt.Errorf("filter databases: %v", err)
	}

	p.logger.With("old-db-count", len(databases)).Info("Databases filtered.")

	return &dropperCloser{
		ctx:       p.ctx,
		logger:    p.logger,
		err:       p.err,
		db:        p.db,
		databases: databases,
	}
}

func postgresFilter(db *sqlx.DB, timestampType internal.TimestampType) *Filter {
	f := &Filter{
		db:    db,
		query: `SELECT datname FROM pg_database WHERE datname ~ $1`,
	}
	switch timestampType {
	case internal.TimestampPrefix:
		f.pattern = PostgresPrefixPattern
		f.getTimeFunc = getPrefixNumPart
	case internal.TimestampPostfix:
		f.pattern = PostgresPostfixPattern
		f.getTimeFunc = getPostfixNumPart
	}

	return f
}
