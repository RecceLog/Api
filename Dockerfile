FROM golang:1.25.5-alpine3.23

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/air-verse/air@latest

COPY . .

RUN go mod download

CMD ["air"]
