# Используем базовый образ для сборки
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем Air
RUN go install github.com/air-verse/air@latest

# Копируем go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o awesome_gin_app .

# Создаем конечный образ
FROM golang:1.23-alpine

# Устанавливаем Air в конечный образ
RUN go install github.com/air-verse/air@latest

# Копируем бинарный файл из предыдущего образа
COPY --from=builder /app/awesome_gin_app /app/awesome_gin_app

WORKDIR /app

# Определяем команду запуска Air
CMD ["air"]
