package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"strings"
	"time"
)

const (
	Postgres = "postgres"
	MsSql    = "mssql"
)

var (
	ttlFlag                                         time.Duration
	dbTypeFlag, connStringFlag, jobScheduleCronFlag string
)

func init() {
	flag.DurationVar(&ttlFlag, "db-ttl", 6*time.Hour, "Database time to live")
	flag.StringVar(&dbTypeFlag, "db-type", Postgres, "DB type. Must be postgres or MsSql")
	flag.StringVar(&connStringFlag, "conn-string", "", "DB connection string")
	flag.StringVar(&jobScheduleCronFlag, "cron", "@every 10s", "Job Schedule in cron format")
	flag.Parse()
}

func main() {
	job, err := getDbDroppingJob(dbTypeFlag)
	if err != nil {
		failWithParsingError(err)
	}
	schedule, err := cron.ParseStandard(jobScheduleCronFlag)
	if err != nil {
		failWithParsingError(err)
	}

	c := cron.New()
	job.Setup(connStringFlag, ttlFlag)
	c.Schedule(schedule, job)
	c.Start()

	fmt.Printf("Start db dropping job for %v with schedule: %v\n", dbTypeFlag, jobScheduleCronFlag)
	fmt.Println(`Type "exit" to stop.`)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := strings.ToLower(scanner.Text())
		if text == "exit" {
			<-c.Stop().Done()
			os.Exit(0)
		}
	}
}

func getDbDroppingJob(dbType string) (DbDroppingJob, error) {
	switch strings.ToLower(dbType) {
	default:
		return nil, fmt.Errorf("invalid db-type: %v", dbType)
	case Postgres:
		return new(PostgresDbDropper), nil
	case MsSql:
		return new(MsSqlDbDropper), nil
	}
}

func failWithParsingError(err error) {
	log.Fatalf("Argument parsing error: %v", err)
}
