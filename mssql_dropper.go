package main

import "fmt"

type MsSqlDbDropper struct {
	DbDropper
}

func (d *MsSqlDbDropper) Run() {
	run(d)
}

func (d *MsSqlDbDropper) GetDbNames() (names []string, err error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (d *MsSqlDbDropper) DropDb(dbName string) error {
	return fmt.Errorf("not implemented")
}
