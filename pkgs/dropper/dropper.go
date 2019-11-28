package dropper

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type DbDropper interface {
	// Find all DBs in collection which satisfy following condition:
	// 'UTC now' subtract 'Db-prefix to unix time' is lesser than DB_TTL
	FindOldDatabases(names []string, ttl time.Duration) <-chan string

	// Drops db with particular name
	DropDb(dbName string) error
}

type baseDbDropper struct{}

func (d *baseDbDropper) FindOldDatabases(names []string, ttl time.Duration) <-chan string {
	ch := make(chan string)
	now := int(time.Now().Unix())
	go func() {
		for _, name := range names {
			utcPrefix := strings.Split(name, "-")[0]
			dbCreatedAt, err := strconv.Atoi(utcPrefix)

			if err != nil {
				log.Printf("Error while parsing dbName's utcPrefix: %v", err)
				continue
			}

			if databaseIsToOld(now, dbCreatedAt, ttl) {
				ch <- name
			}
		}
		close(ch)
	}()
	return ch
}

func (d *baseDbDropper) DropDb(dbName string) error {
	return fmt.Errorf("not implemented")
}

func databaseIsToOld(now int, dbCreatedAt int, ttl time.Duration) bool {
	return now-dbCreatedAt > int(ttl/time.Second)
}
