<<<<<<< HEAD
# Go REST API CRUD с GORM, PostgreSQL и Docker

Проект представляет собой RESTful API на Go с использованием:
- GORM для работы с базой данных
- PostgreSQL в качестве СУБД
- Docker и Docker Compose для контейнеризации
- Gin Web Framework для маршрутизации HTTP-запросов
- Структурированный подход с разделением кода на модули

## Структура проекта

```
go-rest-api/
│
├── cmd/
│   └── api/
│       └── main.go          # Точка входа в приложение
│
├── internal/
│   ├── config/
│   │   └── config.go        # Конфигурация приложения
│   │
│   ├── models/
│   │   └── user.go          # Модели данных (User)
│   │
│   ├── database/
│   │   └── db.go            # Инициализация и управление БД
│   │
│   ├── handlers/
│   │   └── user_handler.go  # Обработчики HTTP-запросов
│   │
│   └── router/
│       └── router.go        # Настройка маршрутов API
│
├── .env                     # Переменные окружения
│
├── go.mod                   # Файл модуля Go
│
├── Dockerfile               # Конфигурация Docker
│
├── docker-compose.yml       # Конфигурация Docker Compose
│
└── README.md                # Документация
```

## Запуск приложения

### С помощью Docker Compose

1. Клонируйте репозиторий
2. Выполните команду:

```bash
docker-compose up -d
```

Это запустит:
- PostgreSQL на порту 5432
- API на порту 8080

### Локальный запуск для разработки

1. Установите зависимости:

```bash
go mod download
```

2. Запустите PostgreSQL (локально или через Docker):

```bash
docker run -d -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres postgres:14-alpine
```

3. Запустите приложение:

```bash
go run cmd/api/main.go
```

## API Endpoints

### Пользователи (Users)

- `GET /api/users` - Получить список всех пользователей
- `GET /api/users/:id` - Получить пользователя по ID
- `POST /api/users` - Создать нового пользователя
- `PUT /api/users/:id` - Обновить пользователя
- `DELETE /api/users/:id` - Удалить пользователя

### Примеры запросов

#### Создание пользователя
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Иван Петров","email":"ivan@example.com","password":"секрет"}'
```

#### Получение всех пользователей
```bash
curl http://localhost:8080/api/users
```

#### Получение пользователя по ID
```bash
curl http://localhost:8080/api/users/1
```

#### Обновление пользователя
```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Иван Сидоров","email":"ivan.new@example.com"}'
```

#### Удаление пользователя
```bash
curl -X DELETE http://localhost:8080/api/users/1
```
=======
# GoRestProject
>>>>>>> 1e5b030465f64e5a08134b2684c416a6064ea493
