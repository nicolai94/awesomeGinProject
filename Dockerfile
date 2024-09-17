FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
# Установите wire
RUN go install github.com/google/wire/cmd/wire@latest

# Сгенерируйте код с помощью wire
RUN wire ./config

RUN go build -o awesome_gin_app .

FROM alpine:latest

COPY --from=builder /app/awesome_gin_app /usr/local/bin/awesome_gin_app

EXPOSE 8080

CMD ["awesome_gin_app"]
