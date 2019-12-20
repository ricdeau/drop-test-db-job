#!/bin/sh

./dropper --db-type "$DB_TYPE" --db-ttl "$DB_TTL" --conn-string "$CONNECTION_STRING" --cron "$JOB_SCHEDULE_CRON"