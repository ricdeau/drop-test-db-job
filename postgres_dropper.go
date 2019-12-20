package main

import _ "github.com/lib/pq"

const (
	postgres               = "postgres"
	postgresQueryDatabases = "SELECT datname FROM pg_database WHERE datname ~ '^\\d{10}-'"
)

// PostgresDbDropper is a dbDropper realisation for pl-sql dialect
type PostgresDbDropper struct {
	dbDropper
}

// Run - runs db dropper
func (d *PostgresDbDropper) Run() {
	run(d)
}

// GetDbNames - gets all database names that satisfies name condition
func (d *PostgresDbDropper) GetDbNames() ([]string, error) {

	return getDbNames(postgres, d.connString, postgresQueryDatabases)
}

// DropDb - drops database with following name
// dbName - name of database to drop
func (d *PostgresDbDropper) DropDb(dbName string) error {
	return dropDb(postgres, dbName, d.connString)
}
