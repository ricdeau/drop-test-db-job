package sql

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ricdeau/drop-test-db-job/internal"
)

type Filter struct {
	db          *sqlx.DB
	query       string
	pattern     string
	getTimeFunc func(string) (time.Time, error)
}

func MySqlFilter(db *sqlx.DB, timestampType internal.TimestampType) *Filter {
	f := &Filter{
		db:    db,
		query: `SELECT schema_name FROM information_schema.schemata WHERE schema_name REGEXP ?`,
	}
	switch timestampType {
	case internal.TimestampPrefix:
		f.pattern = MySqlPrefixPattern
		f.getTimeFunc = getPrefixNumPart
	case internal.TimestampPostfix:
		f.pattern = MySqlPostfixPattern
		f.getTimeFunc = getPostfixNumPart
	}

	return f
}

func MSSQLFilter(db *sqlx.DB, timestampType internal.TimestampType) *Filter {
	f := &Filter{
		db:    db,
		query: `select name from sys.databases where name like '[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]-%'`,
	}
	switch timestampType {
	case internal.TimestampPrefix:
		f.pattern = MSSQLPrefixPattern
		f.getTimeFunc = getPrefixNumPart
	case internal.TimestampPostfix:
		f.pattern = MSSQLPostfixPattern
		f.getTimeFunc = getPostfixNumPart
	}

	return f
}

func (f *Filter) FilterDatabases(ctx context.Context, createdBefore time.Time) ([]string, error) {
	names := []string{}
	err := f.db.SelectContext(ctx, &names, f.query, f.pattern)
	if err != nil {
		return nil, fmt.Errorf("query databases: %v", err)
	}

	result := []string{}
	for _, name := range names {
		ts, err := f.getTimeFunc(name)
		if err != nil {
			return nil, fmt.Errorf("get timestamp from dbName: %q:%v", name, err)
		}

		if ts.Before(createdBefore) {
			result = append(result, name)
		}
	}

	return result, nil
}

func getPrefixNumPart(name string) (t time.Time, err error) {
	return getCreationTime(name[:13])
}

func getPostfixNumPart(name string) (t time.Time, err error) {
	return getCreationTime(name[len(name)-13:])
}

func getCreationTime(timeString string) (t time.Time, err error) {
	millis, err := strconv.Atoi(timeString)
	if err != nil {
		return t, fmt.Errorf("invalid timestamp type")
	}

	return time.UnixMilli(int64(millis)), nil
}
