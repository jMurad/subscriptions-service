# Subscriptions Service

<div align="center">

Production-ready REST API сервис для управления пользовательскими подписками
с поддержкой агрегации стоимости, Swagger-документации, Docker и PostgreSQL.

<br>

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge\&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-316192?style=for-the-badge\&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=for-the-badge\&logo=docker)
![Swagger](https://img.shields.io/badge/Swagger-OpenAPI-85EA2D?style=for-the-badge\&logo=swagger)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

</div>

---

## ✨ Features

* CRUD операции для подписок
* Агрегация общей стоимости подписок
* Фильтрация и пагинация
* Soft Delete
* RESTful API
* PostgreSQL + pgx
* Docker / Docker Compose
* Swagger/OpenAPI документация
* Graceful Shutdown
* Middleware logging/recovery
* Unit и Integration тесты
* Clean Architecture

---

## 🏗 Architecture

```text
.
├── cmd/
│   └── app/                 # Точка входа приложения
│
├── internal/
│   ├── config/              # Конфигурация приложения
│   ├── model/               # Доменные модели и интерфейсы
│   ├── service/             # Бизнес-логика
│   ├── repository/          # Работа с PostgreSQL
│   ├── transport/http/      # HTTP handlers, DTO, routes
│   ├── middleware/          # Middleware
│   └── app/                 # Инициализация приложения
│
├── migrations/              # SQL миграции
├── docs/                    # Swagger документация
├── docker-compose.yml
├── Makefile
└── README.md
```

---

## 🚀 Quick Start

### 1. Clone repository

```bash
git clone https://github.com/jMurad/subscriptions-service.git

cd subscriptions-service
```

---

### 2. Configure environment

Создай `.env` файл:

```env
APP_PORT=8080
DATABASE_URL=postgres://postgres:postgres@postgres:5432/subscriptions?sslmode=disable
```

---

### 3. Run with Docker

```bash
docker compose up --build
```

Приложение будет доступно:

* API → `http://localhost:8080`
* Swagger → `http://localhost:8080/swagger/index.html`

---

## 📚 API Endpoints

| Method | Endpoint                         | Description                 |
| ------ | -------------------------------- | --------------------------- |
| POST   | `/subscriptions`                 | Создать подписку            |
| GET    | `/subscriptions/{id}`            | Получить подписку           |
| GET    | `/subscriptions/user/{user_id}`  | Получить подписку по userID |
| PUT    | `/subscriptions/{id}`            | Обновить подписку           |
| DELETE | `/subscriptions/{id}`            | Удалить подписку            |
| GET    | `/subscriptions`                 | Список подписок             |
| GET    | `/subscriptions/total`           | Общая стоимость             |

---

## 📦 Subscription Model

```json
{
  "id": 1,
  "user_id": 42,
  "service_name": "Netflix",
  "price": 999,
  "start_date": "11-2025",
  "end_date": "12-2025",
  "created_at": "01-2025",
}
```

---

## 🧪 Running Tests

### Unit Tests

```bash
make test-service
make test-handler
make test-middleware
```

### Integration Tests

```bash
make test-repo
```

---

## 🛠 Make Commands

| Command             | Description        |
| ------------------- | ------------------ |
| `make run`          | Запуск приложения  |
| `make build`        | Сборка проекта     |
| `make test`         | Тесты              |
| `make migrate-up`   | Применить миграции |
| `make migrate-down` | Откатить миграции  |
| `make swagger`      | Генерация Swagger  |
| `make docker-up`    | Запуск Docker      |
| `make docker-down`  | Остановка Docker   |

---

## 📖 Swagger

После запуска приложения документация доступна:

```text
http://localhost:8080/swagger/index.html
```

---

## 🧱 Tech Stack

* Go 1.25
* PostgreSQL 17
* pgx
* chi router
* Docker
* Swagger / OpenAPI
* Testify

---

## 🔒 Error Handling

Проект использует централизованную обработку ошибок:

* typed errors
* errors.Is support
* HTTP error mapping
* validation errors
* internal server protection

---

## 📈 Project Goals

Проект создан для демонстрации:

* backend архитектуры
* работы с PostgreSQL
* REST API разработки
* Docker инфраструктуры
* тестирования Go приложений
* production-style подхода к разработке

---

## 📄 License

MIT License

---

<div align="center">

Made with Go ❤️

</div>
