# Auth Service

Небольшой `gRPC`-сервис аутентификации на Go.

Что умеет:
- регистрация пользователя;
- логин с выдачей JWT-токена;
- проверка, является ли пользователь администратором.

## Стек

- Go
- gRPC
- SQLite
- golang-migrate

## Конфиг

Основной конфиг лежит в `config/config.yaml`.

Пример:

```yaml
env: "local"
storage_path: "./storage/sso.db"
token_ttl: 1h

grpc:
  port: 44044
  timeout: 10h
```

## Запуск миграций

```bash
go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations
```

## Запуск сервиса

```bash
go run ./cmd/sso/main.go --config=./config/config.yaml
```

## API

Proto-файл лежит в `proto/sso/sso.proto`.
