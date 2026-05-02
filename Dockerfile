FROM golang:1.25.5-alpine3.23

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.27.0

COPY . .
RUN go mod download

RUN mkdir /app/bin
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd

CMD ["/app/bin/api"]
