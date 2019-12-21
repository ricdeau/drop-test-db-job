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

// Db engine type, for which job is scheduled
const (
	Postgres = "postgres"
	MsSQL    = "mssql"
)

const (
	defaultTtl  = 6 * time.Hour
	defaultCron = "@every 1m"
)

// Application keys
const (
	backgroundKey = "bg"
	ttlKey        = "db-ttl"
	dbTypeKey     = "db-type"
	connStringKey = "conn-string"
	cronKey       = "cron"
)

var (
	background      bool
	dbTtl           time.Duration
	dbType          string
	connString      string
	jobScheduleCron string
)

func init() {
	flag.BoolVar(&background, backgroundKey, false, "Set application to background")
	flag.DurationVar(&dbTtl, ttlKey, defaultTtl, "Database time to live")
	flag.StringVar(&dbType, dbTypeKey, Postgres, "DB type. Must be postgres or MsSQL")
	flag.StringVar(&connString, connStringKey, "", "DB connection string")
	flag.StringVar(&jobScheduleCron, cronKey, defaultCron, "Job Schedule in cron format")
	flag.Parse()
}

func main() {
	job, err := getDbDroppingJob(dbType)
	if err != nil {
		failWithParsingError(err)
	}
	schedule, err := cron.ParseStandard(jobScheduleCron)
	if err != nil {
		failWithParsingError(err)
	}

	c := cron.New()
	job.Setup(connString, dbTtl)
	c.Schedule(schedule, job)
	fmt.Printf("Start db dropping job for %v with schedule: %v\n", dbType, jobScheduleCron)
	if background {
		c.Run()
	} else {
		fmt.Println(`Type "exit" to stop.`)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := strings.ToLower(scanner.Text())
			if text == "exit" {
				stop(c)
			}
		}
	}
}

func stop(c *cron.Cron) {
	fmt.Println("Stopping...")
	<-c.Stop().Done()
	os.Exit(0)
}

func getDbDroppingJob(dbType string) (DbDroppingJob, error) {
	switch strings.ToLower(dbType) {
	default:
		return nil, fmt.Errorf("invalid db-type: %v", dbType)
	case Postgres:
		return new(PostgresDbDropper), nil
	case MsSQL:
		return new(MsSQLDbDropper), nil
	}
}

func failWithParsingError(err error) {
	log.Fatalf("Argument parsing error: %v", err)
}
