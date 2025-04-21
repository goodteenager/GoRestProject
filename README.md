# Структура проекта Go REST API

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
│   ├── migrations/         # Новый пакет для миграций
│   │   ├── migrations.go   # Основной код системы миграций
│   │   └── 001_create_users_table.go  # Первая миграция
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

# Go REST API CRUD с GORM, PostgreSQL, миграциями и Docker

Проект представляет собой RESTful API на Go с использованием:
- GORM для работы с базой данных
- PostgreSQL в качестве СУБД
- Система миграций для управления схемой БД
- Docker и Docker Compose для контейнеризации
- Gin Web Framework для маршрутизации HTTP-запросов
- Структурированный подход с разделением кода на модули

## Основные компоненты

### 1. Система миграций
В проект добавлена система миграций, которая:
- Отслеживает версии схемы базы данных
- Автоматически применяет только новые миграции
- Выполняет миграции в транзакционном режиме
- Хранит историю выполненных миграций в таблице `migrations`

### 2. REST API для работы с пользователями
- Полный CRUD функционал (Create, Read, Update, Delete)
- Валидация входных данных
- Безопасное хранение данных пользователей

### 3. Конфигурация и окружение
- Поддержка переменных окружения через файл `.env`
- Гибкая конфигурация компонентов приложения

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

## Система миграций

### Как создать новую миграцию

1. Создайте новый файл в директории `internal/migrations/` с префиксом, указывающим порядок выполнения:

```go
// internal/migrations/002_add_role_to_users.go

package migrations

import (
    "gorm.io/gorm"
)

func init() {
    Register("002_add_role_to_users", addRoleToUsers)
}

func addRoleToUsers(db *gorm.DB) error {
    return db.Exec(`
        ALTER TABLE users ADD COLUMN role VARCHAR(50) DEFAULT 'user' NOT NULL;
    `).Error
}
```

2. При следующем запуске приложения миграция будет автоматически применена.

## Преимущества данной архитектуры

1. **Модульность и масштабируемость**:
    - Четкое разделение ответственности между компонентами
    - Легкое добавление новых функций и модулей

2. **Контроль версий базы данных**:
    - Отслеживание изменений схемы БД
    - Безопасные обновления в производственной среде

3. **Контейнеризация**:
    - Изолированное окружение для разработки и деплоя
    - Простое развертывание на разных платформах