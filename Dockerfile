FROM golang:1.23-alpine AS builder

# Установка рабочей директории
WORKDIR /app

# Копирование и загрузка зависимостей
COPY go.mod ./
RUN go mod download

# Копирование исходного кода
COPY . ./

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Финальный образ
FROM alpine:latest

# Установка сертификатов для HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копирование бинарного файла из этапа сборки
COPY --from=builder /app/api ./api
COPY --from=builder /app/.env .

# Открытие порта приложения
EXPOSE 8080

# Запуск приложения
CMD ["./api"]