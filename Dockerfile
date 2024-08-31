# Build app
FROM golang:1.21.8-alpine3.19 AS builder

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN apk add --no-cache gcc musl-dev
RUN go mod download

COPY ./cmd/ ./cmd
COPY ./config ./config
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./web/ ./web
COPY ./*.db .

RUN go build -ldflags='-s -w -extldflags "-static"' -o /scheduler ./cmd/app/main.go

# Container with app
FROM alpine

ENV TODO_PORT=7540 \
    TODO_DBFILE=scheduler.db \
    TODO_PASSWORD=123

WORKDIR /app

COPY --from=builder /scheduler .
COPY --from=builder /app/*.db .
COPY --from=builder /app/web ./web

EXPOSE ${TODO_PORT}

CMD ["/app/scheduler"]