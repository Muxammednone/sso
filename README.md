# SSO — сервис авторизации на Go

Учебный gRPC-сервис регистрации и входа. Хранит аккаунты в PostgreSQL, хеширует пароли через bcrypt и выдаёт JWT. Используется в связке с [TodoApp API](https://github.com/Muxammednone/golang-todoapp).

## Возможности и устройство

- `Register`: создание аккаунта по email и паролю.
- `Login`: проверка пароля и выдача JWT HS256 с `sub`, `iat`, `exp`.
- PostgreSQL через pgxpool и SQL-миграции.
- Настройки через переменные окружения, таймауты и корректное завершение сервера.
- Ошибки приложения преобразуются в gRPC status codes.

```text
TodoApp / gRPC-клиент → gRPC transport → service → PostgreSQL repository
                                          |
                                          └── JWT
```

Код: `cmd/sso` — запуск; `internal/config` — настройки; `internal/transport/grpc` — обработчики; `internal/service` — регистрация и вход; `internal/repository/postgres` — запросы; `internal/token` — JWT; `migrations` — схема базы.

Контракт находится в [sso-proto](https://github.com/Muxammednone/sso-proto), файл `proto/sso/sso.proto`. Сгенерированный Go-код подключён как модуль `v0.4.0`: для обычного запуска генерация не требуется.

## Требования

Go 1.26.4 или новее (для всей связки с TodoApp — 1.26.5+), Docker с Compose и Task CLI. На Windows нужен запущенный Docker Desktop с Linux containers.

## Запуск локально

Скопируйте `.env.example` в `.env` и заполните значения:

```dotenv
POSTGRES_USER=sso
POSTGRES_PASSWORD=sso_local_password
POSTGRES_DB=sso_local_dev
DATABASE_URL=postgres://sso:sso_local_password@127.0.0.1:5433/sso_local_dev?sslmode=disable
JWT_SECRET=<случайный секрет длиной не меньше 32 байт>
GRPC_ADDR=127.0.0.1:50051
TOKEN_TTL=1h
DB_TIMEOUT=5s
SHUTDOWN_TIMEOUT=5s
```

Пароль базы в примере предназначен для локальной разработки. Замените `JWT_SECRET` собственным значением и используйте тот же секрет в TodoApp. `.env` не должен попадать в Git.

Сгенерировать секрет в PowerShell можно так:

```powershell
$bytes = New-Object byte[] 32
$rng = [System.Security.Cryptography.RandomNumberGenerator]::Create()
$rng.GetBytes($bytes)
[Convert]::ToBase64String($bytes)
$rng.Dispose()
```

Из корня репозитория:

```sh
task db:up
task db:migrate
task run
```

Если миграция получила ошибку соединения сразу после запуска базы, дождитесь её готовности и повторите `task db:migrate`.

PostgreSQL доступна с хоста по `127.0.0.1:5433`; внутри Compose мигратор обращается к `postgres:5432`. gRPC-сервер слушает `127.0.0.1:50051`.

| Команда | Назначение |
| --- | --- |
| `task run` | Запуск сервиса с загрузкой `.env` |
| `task build` | Сборка Go-пакетов |
| `task db:up` | Запуск PostgreSQL |
| `task db:migrate` | Применение миграций |
| `task db:status` | Состояние контейнеров |
| `task db:down` | Остановка контейнеров с сохранением тома базы |

Прямой `go run ./cmd/sso` не загружает `.env`: переменные должны уже находиться в окружении.

## Конфигурация

| Переменная | Значение |
| --- | --- |
| `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` | Настройки контейнера PostgreSQL и мигратора |
| `DATABASE_URL` | Обязательная строка подключения приложения |
| `JWT_SECRET` | Обязательный общий секрет HS256, минимум 32 байта |
| `GRPC_ADDR` | Адрес сервера, по умолчанию `127.0.0.1:50051` |
| `TOKEN_TTL` | Срок жизни JWT, по умолчанию `1h` |
| `DB_TIMEOUT` | Таймаут работы с базой, по умолчанию `5s` |
| `SHUTDOWN_TIMEOUT` | Таймаут завершения, по умолчанию `5s` |

Длительности должны быть положительными. Для использования значения по умолчанию удалите переменную из `.env`, а не оставляйте пустую строку.

## Проверка через Postman gRPC

Создайте gRPC-запрос к `127.0.0.1:50051` без TLS и импортируйте `proto/sso/sso.proto` из репозитория контрактов.

Вызовите `sso.SSO/Register`:

```json
{"email":"demo@example.com","password":"Demo-password-123!"}
```

Ответ содержит `user_id`. Затем вызовите `sso.SSO/Login` с тем же телом. Ответ содержит `access_token`. В HTTP-запросах TodoApp передавайте его как `Authorization: Bearer <access_token>`.

SSO создаёт аккаунт, но не профиль TodoApp. После входа профиль создаётся отдельно через `POST /api/v1/me` в TodoApp. Полный HTTP-сценарий приведён в [README TodoApp](https://github.com/Muxammednone/golang-todoapp#readme).

## Границы проекта

Сервис реализует регистрацию и вход. Refresh-токены, отзыв JWT, восстановление пароля и подтверждение email не реализованы. Локальный gRPC работает без TLS; это не готовая конфигурация публичного развёртывания.
