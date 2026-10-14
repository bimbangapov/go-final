FROM golang:1.26.2-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todo-app .

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

CMD ["./todo-app"]