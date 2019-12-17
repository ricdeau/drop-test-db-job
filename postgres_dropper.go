package main

// PostgresDbDropper is a DbDropper realisation for pl-sql dialect
type PostgresDbDropper struct {
	DbDropper
}

// Run - runs db dropper
func (d *PostgresDbDropper) Run() {
	run(d)
}

// GetDbNames - gets all database names that satisfies name condition
func (d *PostgresDbDropper) GetDbNames() ([]string, error) {
	return getDbNames("postgres", d.connString, "SELECT datname FROM pg_database WHERE datname ~ '^\\d{10}-'")
}

// DropDb - drops database with following name
// dbName - name of database to drop
func (d *PostgresDbDropper) DropDb(dbName string) error {
	return dropDb("postgres", dbName, d.connString)
}
