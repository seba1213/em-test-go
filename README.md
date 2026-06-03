# em-test-go
Тестовое задание Junior Golang Developer Effective Mobile

## Dev Container

Окружение для разработки описано в `.devcontainer/` и поднимает два сервиса:

| Сервис | Назначение |
|--------|------------|
| `app` | Go 1.25, рабочая директория `/workspaces/em-test-go` |
| `db` | PostgreSQL |

Переменные окружения задаются в `.devcontainer/.env` и подключаются к обоим сервисам.

### Запуск окружения

1. Откройте репозиторий в VS Code или Cursor.
2. Выполните команду **Dev Containers: Reopen in Container**.
3. Дождитесь сборки контейнера и установки зависимостей.

Порт `8080` пробрасывается на хост — API доступен с локальной машины по адресу `http://localhost:8080`.

### Запуск приложения

```bash
go run main.go
```

Для заполнения базы демо-данными:

```bash
go run cmd/demo_filler/main.go
```

## Production

Production-окружение описано в `deploy/` и поднимает два сервиса через Docker Compose:

| Сервис | Назначение |
|--------|------------|
| `app` | API-сервер (образ собирается из `Dockerfile` в корне репозитория) |
| `db` | PostgreSQL 16.3 |


### Подготовка

```bash
cd deploy
cp .env.example .env
```

### Запуск

```bash
cd deploy
docker compose up -d --build
```