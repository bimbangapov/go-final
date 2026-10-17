FROM golang:1.26.2-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app .

FROM alpine:latest as application
WORKDIR /srv
COPY --from=builder /app/todo-app ./todo-app
COPY --from=builder /app/web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

CMD ["./todo-app"]