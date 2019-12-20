package main

import _ "github.com/denisenkom/go-mssqldb"

const (
	sqlserver           = "sqlserver"
	mssqlQueryDatabases = "select name from sys.databases where name like '[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]-%'"
)

// MsSQLDbDropper is a dbDropper realisation for T-SQL dialect
type MsSQLDbDropper struct {
	dbDropper
}

// Run - runs db dropper
func (d *MsSQLDbDropper) Run() {
	run(d)
}

// GetDbNames - gets all database names that satisfies name condition
func (d *MsSQLDbDropper) GetDbNames() (names []string, err error) {
	return getDbNames(sqlserver, d.connString, mssqlQueryDatabases)
}

// DropDb - drops database with following name
// dbName - name of database to drop
func (d *MsSQLDbDropper) DropDb(dbName string) error {
	return dropDb(sqlserver, dbName, d.connString)
}
