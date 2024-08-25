# Build app
FROM golang:1.22 AS builder

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY *.db ./
COPY ./web/ ./web/

RUN go build -o /scheduler

# Container with app
FROM ubuntu

ENV TODO_PORT=7540 \
    TODO_DBFILE=scheduler.db

WORKDIR /app

COPY --from=builder /scheduler .
COPY --from=builder /app/*.db .
COPY --from=builder /app/web ./web

EXPOSE ${TODO_PORT}

CMD ["/app/scheduler"]