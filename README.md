# Основной веб-сервис SQL Analyzer

## Краткое описание

Основной backend учебной системы SQL Analyzer для курса «Разработка интернет-приложений». Сервис предоставляет API для работы с индексами, заявками на расчёт времени выполнения SQL-запросов и пользователями.

## Основные возможности

- получение, создание, изменение и удаление индексов;
- добавление индексов в заявку и изменение связей «индекс — заявка»;
- создание, просмотр, изменение, формирование, удаление и завершение заявок;
- регистрация, вход, выход и просмотр профиля пользователя;
- разграничение доступа для пользователя и модератора;
- передача результатов асинхронного расчёта через endpoint обновления результата;
- Swagger-документация в [`docs/swagger.yaml`](docs/swagger.yaml).

## Лабораторные работы и ветки

- [Лабораторная работа № 1 — хранение данных в памяти](https://github.com/colan1ch/sql-analyzer-web-service/tree/sql_analyzer_srr_in_memory)
- [Лабораторная работа № 2 — PostgreSQL и ORM](https://github.com/colan1ch/sql-analyzer-web-service/tree/sql_analyzer_db_implementation)
- [Лабораторная работа № 3 — веб-сервис для SPA](https://github.com/colan1ch/sql-analyzer-web-service/tree/SPA_backend)
- [Лабораторная работа № 4 — авторизация и Swagger](https://github.com/colan1ch/sql-analyzer-web-service/tree/SPA_backend_swagger)

В `main` находится финальная версия ветки `SPA_backend_swagger`.

## Стек технологий

- Go `1.24.2`;
- Gin;
- GORM и PostgreSQL driver;
- PostgreSQL;
- Redis `6.2` в Docker Compose;
- MinIO;
- JWT-аутентификация и middleware проверки ролей;
- Swagger UI и swagger-файлы;
- Docker Compose для PostgreSQL, Redis, MinIO, Adminer и Nginx.

## Установка и запуск

Требуются Go, Docker и Docker Compose. Параметры приложения находятся в [`config/config.toml`](config/config.toml), параметры подключения к внешним сервисам — в `.env`.

Запуск инфраструктуры:

```bash
docker compose up -d
```

Применение миграций:

```bash
go run ./cmd/migrate
```

Запуск веб-сервиса:

```bash
go run ./cmd/LAB1
```

Сервис слушает адрес `0.0.0.0:8080` согласно [`config/config.toml`](config/config.toml). Файл `docker-compose.yml` запускает PostgreSQL на `5432`, Adminer на `8081`, Redis на `6379`, MinIO API на `9000` и консоль MinIO на `9001`.

## API и важные файлы

Swagger-описание доступно в [`docs/swagger.yaml`](docs/swagger.yaml), [`docs/swagger.json`](docs/swagger.json) и [`docs/docs.go`](docs/docs.go). Маршруты и обработчики находятся в [`internal/app/handler`](internal/app/handler), доступ к данным — в [`internal/app/repository`](internal/app/repository), доменные структуры — в [`internal/app/ds`](internal/app/ds).

Ключевые группы API:

```text
/api/v1/indexes
/api/v1/queries
/api/v1/indexes_query
/api/v1/users
```

Точные HTTP-методы, параметры и схемы ответов следует смотреть в Swagger, чтобы не дублировать спецификацию в README.

## Структура проекта

```text
cmd/LAB1/                 точка запуска сервиса
cmd/migrate/              миграция схемы PostgreSQL через GORM
config/                   конфигурация адреса сервиса
docs/                     Swagger-документация
internal/app/api_types/   структуры JSON API
internal/app/ds/          доменные структуры
internal/app/handler/     HTTP-обработчики и middleware
internal/app/repository/ доступ к PostgreSQL
internal/app/minioClient/ клиент MinIO
internal/pkg/             запуск приложения
resources/                изображения и CSS
templates/                HTML-шаблоны
docker-compose.yml        инфраструктура разработки
nginx.conf                конфигурация Nginx для MinIO
```

## Статус проекта

Учебный backend курса «Разработка интернет-приложений». В `main` объединена финальная версия `SPA_backend_swagger`; асинхронное вычисление вынесено в отдельный репозиторий.

## Связанные репозитории

- [Фронтенд](https://github.com/colan1ch/sql-analyzer)
- [Асинхронный сервис](https://github.com/colan1ch/sql-analyzer-async-web-service)

## Описание для GitHub

Основной веб-сервис SQL Analyzer для заявок на расчёт производительности SQL-запросов. Стек: Go, Gin, GORM, PostgreSQL, Redis, MinIO, Docker Compose и Swagger. Реализует API индексов, заявок, пользователей и интеграцию с асинхронным расчётом.
