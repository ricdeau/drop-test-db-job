package dropper

import "time"

type DbDropper interface {
	// Find all DBs in collection which satisfy following condition:
	// 'UTC now' subtract 'Db-prefix to unix time' is lesser than DB_TTL
	FindOldDatabases(names []string, ttl time.Duration) <-chan string

	// Drops db with particular name
	DropDb(dbName string) error
}
