package main

import (
	"database/sql"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// DbDroppingJob - seek for old databases than drops them.
// The name of database must begin with 10 decimal digits, which represents
// UNIX milliseconds timestamp of database creation time
type DbDroppingJob interface {
	cron.Job
	Setup(connString string, ttl time.Duration)
}

type dbDropperInterface interface {
	GetDbNames() ([]string, error)
	DropDb(dbName string) error
	FilterOldDatabases(names []string) <-chan string
}

type dbDropper struct {
	connString string
	ttl        time.Duration
}

func (d *dbDropper) Setup(connString string, ttl time.Duration) {
	d.connString = connString
	d.ttl = ttl
}

func (d *dbDropper) FilterOldDatabases(names []string) <-chan string {
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

			if databaseIsToOld(now, dbCreatedAt, d.ttl) {
				ch <- name
			}
		}
		close(ch)
	}()
	return ch
}

func run(d dbDropperInterface) {
	log.Println("Start executing job.")
	log.Println("Scanning for databases...")
	names, err := d.GetDbNames()
	if err != nil {
		failWithJobError(err)
	}
	log.Printf("%d databases found\n", len(names))
	log.Println("Start seeking and dropping old databases...")

	var dropped int32
	sema := WaitingSemaphore{counter: 4}
	for name := range d.FilterOldDatabases(names) {
		sema.Acquire()
		go func(n string) {
			defer sema.Release()
			err := d.DropDb(n)
			if err != nil {
				log.Printf("Job execution error: error while dropping %s: %v\n", n, err)
				return
			}
			atomic.AddInt32(&dropped, 1)
		}(name)
	}
	sema.Wait()

	log.Printf("%d old databases dropped\n", dropped)
	log.Println("Job executed")
}

func getDbNames(driverName, connString, query string) (names []string, err error) {
	db, err := sql.Open(driverName, connString)
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return
		}
		names = append(names, name)
	}
	return
}

func dropDb(dbType, dbName, connString string) error {
	db, err := sql.Open(dbType, connString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`DROP DATABASE "%s";`, dbName)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	log.Printf("Database %s has been dropped\n", dbName)
	return nil
}

func databaseIsToOld(now int, dbCreatedAt int, ttl time.Duration) bool {
	return now-dbCreatedAt > int(ttl/time.Second)
}

func failWithJobError(err error) {
	log.Fatalf("Job execution error: %v", err)
}
