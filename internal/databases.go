package internal

import (
	"fmt"
	"strings"
)

type DatabaseType byte

const (
	_ DatabaseType = iota
	DatabasePostgres
	DatabaseMySQL
	DatabaseMSSQL
)

func (d *DatabaseType) UnmarshalJSON(bytes []byte) error {
	var v DatabaseType

	s := strings.Trim(string(bytes), `"`)
	switch strings.ToLower(s) {
	case "postgres", "postgresql":
		v = DatabasePostgres
	case "mysql":
		v = DatabaseMySQL
	case "mssql", "sqlserver":
		v = DatabaseMSSQL
	default:
		return fmt.Errorf("unknown DatabaseType: %q", s)
	}

	*d = v
	return nil
}

func (d DatabaseType) DriverName() string {
	switch d {
	case DatabasePostgres:
		return "postgres"
	case DatabaseMySQL:
		return "mysql"
	case DatabaseMSSQL:
		return "sqlserver"
	default:
		return ""
	}
}
