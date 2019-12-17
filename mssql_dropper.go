package main

import "fmt"

// MsSQLDbDropper is a DbDropper realisation for T-SQL dialect
type MsSQLDbDropper struct {
	DbDropper
}

// Run - runs db dropper
func (d *MsSQLDbDropper) Run() {
	run(d)
}

// GetDbNames - gets all database names that satisfies name condition
func (d *MsSQLDbDropper) GetDbNames() (names []string, err error) {
	return nil, fmt.Errorf("not implemented yet")
}

// DropDb - drops database with following name
// dbName - name of database to drop
func (d *MsSQLDbDropper) DropDb(dbName string) error {
	return fmt.Errorf("not implemented")
}
