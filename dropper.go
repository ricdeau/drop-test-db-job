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

type DbDroppingJob interface {
	cron.Job
	Setup(connString string, ttl time.Duration)
}

type dbDropper interface {
	GetDbNames() ([]string, error)
	DropDb(dbName string) error
	FilterOldDatabases(names []string) <-chan string
}

type DbDropper struct {
	connString string
	ttl        time.Duration
}

func (d *DbDropper) Setup(connString string, ttl time.Duration) {
	d.connString = connString
	d.ttl = ttl
}

func (d *DbDropper) FilterOldDatabases(names []string) <-chan string {
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

func run(d dbDropper) {
	log.Println("Start executing job.")
	log.Println("Scanning for databases...")
	names, err := d.GetDbNames()
	if err != nil {
		failWithJobError(err)
	}
	log.Printf("%d databases found\n", len(names))
	log.Println("Start seeking and dropping old databases...")

	var dropped int32
	sem := WaitingSemaphore{counter: 4}
	for name := range d.FilterOldDatabases(names) {
		sem.Acquire()
		go func(n string) {
			defer sem.Release()
			err := d.DropDb(n)
			if err != nil {
				log.Printf("Job execution error: error while dropping %s: %v\n", n, err)
				return
			}
			atomic.AddInt32(&dropped, 1)
		}(name)
	}
	sem.Wait()

	log.Printf("%d old databases dropped\n", dropped)
	log.Println("Job executed")
}

func getDbNames(dbName, connString, query string) (names []string, err error) {
	db, err := sql.Open(dbName, connString)
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
