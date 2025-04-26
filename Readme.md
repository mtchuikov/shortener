
# shortener

*shortener* - сервис для создания коротких URL.

## Загрузка

* Загрузите последнюю версию бинарного файла со страницы релизов. Для этого выполните команду:
```
TODO
```
* Или загрузите Docker-образ из репозитория:
```
docker pull ghcr.io/mtchuikov/shortener:latest
```
* Или загрузите исходный код и скомпилируйте его самостоятельно:
```
git clone https://github.com/mtchuikov/shortener
sqlc generate
go build -ldflags="-s -w" -o ./build/shortener ./cmd/main.go
```

## Локальный запуск

Для локального запуска следует использовать команду `docker compose up`. Она соберет бинарный файл приложения, а также запустит базу PostgreSQL в Docker-контейнерах.

## Конфигурация

Приложение может быть сконфигурировано при помощи переменных окружения или флагов командной строки (переменные окружения имеют приоритет над флагами):

| Flag                 | Shorthand | Default value               | Description                                                       | Environment Variable    |
|----------------------|-----------|-----------------------------|-------------------------------------------------------------------|-------------------------|
| `--log.file.enable`  | *(none)*          | `false`                     | Switch logging output to a file                                   | `LOG_TO_FILE`           |
| `--log.file`         | *(none)*          | `shortener.log`             | File path where logs will be written                              | `LOG_FILE`              |
| `--log.file.level`   | *(none)*          | `error`                     | Minimum log level to write into the file                          | `LOG_FILE_LEVEL`        |
| `--addr`             | `-a`      | `127.0.0.1:8080`            | Address and port where the HTTP server listens                    | `SERVER_ADDRESS`        |
| `--base`             | `-b`      | `http://127.0.0.1:8080/`    | Base URL used to construct shortened links (base + random ID)     | `BASE_URL`              |
| `--secret`           | `-s`      | `jwtsecret`                 | Secret key for signing and verifying JWTs                         | `JWT_SECRET`            |
| `--dsn`              | `-d`      | *(none)*                    | DSN (Data Source Name) for connecting to the PostgreSQL database  | `DATABASE_DSN`          |
| `--file`             | `-f`      | `shortener.backup`          | File path for backing up shortened URLs from in-memory cache      | `FILE_STORAGE_PATH`     |
| `--verbose`          | `-v`      | `false`                     | Enable verbose (debug-level) logging                              | `VERBOSE`               |
| `--help`             | `-h`      | *(none)*                    | Show help information about supported flags and usage             | *(none)*               |
