FROM golang:alpine AS backend-builder
RUN apk update \
    && apk upgrade \
    && apk add --no-cache git
WORKDIR .
RUN mkdir ./out
RUN GOOS=linux go build -o /out/dropper

FROM alpine:latest
WORKDIR /app
COPY --from=backend-builder /out .
CMD ["./dropper --db-type=$DBTYPE --db-ttl=$DB_TTL --conn-string=$CONNECTION_STRING --cron=$JOB_SCHEDULE_CRON"]
