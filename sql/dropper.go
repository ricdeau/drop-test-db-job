package sql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ricdeau/drop-test-db-job/internal"
	"go.uber.org/zap"
)

type dropperCloser struct {
	ctx    context.Context
	logger *zap.SugaredLogger
	err    error
	db     *sqlx.DB

	databases []string
}

func (d *dropperCloser) Drop() internal.Closer {
	if d.err != nil {
		return d
	}

	for _, dbName := range d.databases {
		_, err := d.db.ExecContext(d.ctx, fmt.Sprintf(`DROP DATABSE "%s"`, dbName))
		if err != nil {
			d.err = fmt.Errorf("drop %q: %v", dbName, err)
			break
		}

		d.logger.With("db-name", dbName).Info("Database dropped.")
	}

	return d
}

func (d *dropperCloser) Close() error {
	err := d.db.Close()
	if err != nil {
		err = fmt.Errorf("close connection: %v", err)
		if d.err != nil {
			d.err = fmt.Errorf("%v + %v", err, err)
		} else {
			d.err = err
		}
	}

	return d.err
}
