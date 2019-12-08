package main

type PostgresDbDropper struct {
	DbDropper
}

func (d *PostgresDbDropper) Run() {
	run(d)
}

func (d *PostgresDbDropper) GetDbNames() ([]string, error) {
	return getDbNames("postgres", d.connString, "SELECT datname FROM pg_database WHERE datname ~ '^\\d{10}-'")
}

func (d *PostgresDbDropper) DropDb(dbName string) error {
	return dropDb("postgres", dbName, d.connString)
}
