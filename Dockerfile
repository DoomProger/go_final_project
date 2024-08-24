FROM golang:1.22 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY *.db ./

RUN go build -o /scheduler

# Container with app
FROM alpine:latest

ENV TODO_PORT=7540 \
    TODO_DBFILE=scheduler.db

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /scheduler .
COPY --from=builder /app/*.db .

EXPOSE ${TODO_PORT}

CMD ["/app/scheduler"]