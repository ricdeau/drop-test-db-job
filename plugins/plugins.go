package plugins

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var plugins = map[string]*Plugin{
	"postgres": {
		name:           "postgres",
		driverName:     "postgres",
		filterDbsQuery: `SELECT datname FROM pg_database WHERE datname ~ $1`,
		dropDbQuery:    `DROP DATABASE "%s"`,
	},
	"mssql": {
		name:           "mssql",
		driverName:     "sqlserver",
		filterDbsQuery: `SELECT name FROM sys.databases WHERE name ILIKE ?`,
		dropDbQuery:    `DROP DATABASE "%s"`,
	},
}

type Result struct {
	DbName string
	Err    error
}

const (
	PrefixTimestamp  = "prefix"
	PostfixTimestamp = "postfix"
)

func RegisterPlugin(name string, plugin *Plugin) error {
	if _, ok := plugins[name]; ok {
		return fmt.Errorf("plugin with name=%s already registered", name)
	}

	plugins[name] = plugin

	return nil
}

func GetPlugin(name string) (*Plugin, error) {
	p, ok := plugins[name]
	if !ok {
		return nil, fmt.Errorf("plugin with name=%s not found", name)
	}

	return p, nil
}

type Plugin struct {
	db *sql.DB

	name           string
	driverName     string
	filterDbsQuery string
	dropDbQuery    string
}

func NewPlugin(name string, driverName string, filterDbsQuery string, dropDbQuery string) *Plugin {
	return &Plugin{
		name:           name,
		driverName:     driverName,
		filterDbsQuery: filterDbsQuery,
		dropDbQuery:    dropDbQuery,
	}
}

func (p *Plugin) GetName() string {
	return p.name
}

func (p *Plugin) Connect(dsn string) (err error) {
	p.db, err = sql.Open(p.driverName, dsn)
	if err != nil {
		return err
	}

	return nil
}

func (p *Plugin) Close() error {
	return p.db.Close()
}

func (p *Plugin) DropDb(ctx context.Context, dbName string) error {
	_, err := p.db.ExecContext(ctx, fmt.Sprintf(p.dropDbQuery, dbName))
	return err
}

func (p *Plugin) GetOldDbs(ctx context.Context, namePattern, timestampPlace string, maxAge time.Duration) (<-chan *Result, error) {
	r, err := regexp.Compile(namePattern)
	if err != nil {
		return nil, err
	}

	results := make(chan *Result)
	rows, err := p.db.QueryContext(ctx, p.filterDbsQuery, namePattern)
	if err != nil {
		return nil, err
	}
	go func() {
		defer rows.Close()
		defer close(results)
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				results <- &Result{Err: err}
				continue
			}

			if !r.MatchString(name) {
				continue
			}

			createdAt, err := getCreatedAt(name, timestampPlace)
			if err != nil {
				results <- &Result{Err: err}
				continue
			}

			if time.Since(createdAt) > maxAge {
				results <- &Result{DbName: name}
			}
		}
	}()

	return results, nil
}

func getCreatedAt(dbName, place string) (time.Time, error) {
	if len(dbName) < 13 {
		return time.Time{}, fmt.Errorf("db name must have at least 13 digits (unix milli), len = %d", len(dbName))
	}

	var dur string
	if place == PrefixTimestamp {
		dur = dbName[:13]
	} else if place == PostfixTimestamp {
		dur = dbName[len(dbName)-13:]
	}

	result, err := strconv.Atoi(dur)
	if err != nil {
		return time.Time{}, fmt.Errorf("convert duration: %v", err)
	}

	return time.UnixMilli(int64(result)), nil
}
