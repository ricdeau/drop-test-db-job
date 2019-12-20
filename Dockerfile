FROM golang:alpine AS backend-builder
RUN apk update \
    && apk upgrade \
    && apk add --no-cache git
WORKDIR $GOPATH/src/dropper
COPY . .
RUN GOOS=linux go build -o /out/dropper
COPY run.sh /out/

FROM alpine:latest
WORKDIR /app
COPY --from=backend-builder /out .
RUN chmod 777 run.sh
CMD ["sh", "./run.sh"]
