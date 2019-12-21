FROM golang:alpine AS backend-builder
RUN apk update \
    && apk upgrade \
    && apk add --no-cache git
WORKDIR $GOPATH/src/dropper
COPY . .
RUN GOOS=linux go build -o /out/dropper

FROM alpine:latest
WORKDIR /app
COPY --from=backend-builder /out/dropper /usr/local/bin/
COPY run.sh .
RUN chmod 777 ./run.sh
CMD ["sh", "run.sh"]
