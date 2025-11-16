# Pull Requests Reviewers Service

## Description

A microservice for automatic assignment of reviewers to Pull Requests (PR) within teams, as well as management of teams and users. Interaction is via HTTP API. The service is implemented in Go and uses PostgreSQL for data storage.

## Features
- Automatic assignment of up to two active reviewers from the author's team (excluding the author) to a PR;
- Reassignment of a reviewer to another active team member;
- Retrieve the list of PRs assigned to a specific user;
- Manage teams and user activity;
- Reviewers cannot be changed after a PR is merged.

## How to use

Before starting, you must manually create a PostgreSQL database (for example, using `createdb` or via a database UI). The service will not create the database automatically, only tables.

You need to create a `.env` file with the following content:

```
DATABASE_URL=postgres://user:password@localhost/dbname
```

Before starting, make sure PostgreSQL is accessible from outside localhost. This setup may differ depending on your OS.

Then start the container using Docker Compose:

```bash
docker-compose up --build
```

The service will be available on port `8080`.

### HTTP Request Examples

The `http/` directory contains sample requests for all main endpoints:
- `post_team_add.http` — add a team
- `post_pull_request_create.http` — create a PR
- `post_pull_request_merge.http` — merge a PR
- `post_pull_request_reassign.http` — reassign a reviewer
- `get_team_get.http` — get team members
- `get_users_get_review.http` — get PRs where the user is a reviewer
- `post_user_set_is_active.http` — change user activity

## Project Structure

- `cmd/server/main.go` — entry point
- `internal/api/` — OpenAPI spec, generated types
- `internal/db/` — database logic
- `internal/handler/` — HTTP handlers
- `http/` — HTTP request examples

---

# Pull Requests Reviewers Service

## Описание

Микросервис для автоматического назначения ревьюверов на Pull Request'ы (PR) внутри команд, а также управления командами и пользователями. Взаимодействие происходит через HTTP API. Сервис реализован на Go, используется PostgreSQL для хранения данных.

## Возможности
- Автоматическое назначение до двух активных ревьюверов из команды автора PR (исключая самого автора);
- Переназначение ревьювера на другого активного участника команды;
- Получение списка PR, назначенных конкретному пользователю;
- Управление командами и активностью пользователей;
- Запрет изменения ревьюверов после merge PR.

## Как пользоваться

Перед запуском необходимо вручную создать базу данных PostgreSQL (например, с помощью `createdb` или через UI). Сервис сам создаёт только таблицы, но не саму базу данных.

Необходимо создать файл `.env` со следующим содержимым:

```
DATABASE_URL=postgres://user:password@localhost/dbname
```

Перед запуском необходимо убедиться, что к PostgreSQL есть доступ из-под неlocalhost. Для каждой операционной системы это настраивается по-разному :(

Далее необходимо запустить контейнер при помощи Docker Compose:

```bash
docker-compose up --build
```

Сервис будет доступен на порту `8080`.

### Примеры HTTP-запросов

В директории `http/` приведены примеры запросов для всех основных эндпоинтов:
- `post_team_add.http` — добавление команды
- `post_pull_request_create.http` — создание PR
- `post_pull_request_merge.http` — merge PR
- `post_pull_request_reassign.http` — переназначение ревьювера
- `get_team_get.http` — получить состав команды
- `get_users_get_review.http` — получить PR'ы, где пользователь назначен ревьювером
- `post_user_set_is_active.http` — смена активности пользователя

## Структура проекта

- `cmd/server/main.go` — точка входа
- `internal/api/` — OpenAPI спецификация, автогенерированные типы
- `internal/db/` — работа с БД
- `internal/handler/` — HTTP-обработчики
- `http/` — примеры HTTP-запросов
