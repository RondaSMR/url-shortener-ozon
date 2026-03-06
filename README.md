# Сокращатель ссылок

---

## Эндопоинты

Сервер принимает два запроса:

### POST / - создание сокращенной ссылки

Ожидаемое тело запроса:
```json
  {
  "url": "https://ozon.ru"
  }
```

Возможные ответы:
* `201 Created` - сокращенная ссылка успешно создана
```json
  {
  "success": true,
  "message": "",
  "result": {
  "url": "i9F9AvSymi"
  }
  }
```
* `400 Bad Request` - невалидный JSON или некорректный формат ссылки (допускаются только http:// и https://)
```json
  {
  "message": "validation error",
  "code": "PS-00102"
  }
```
* `500 Internal Server Error` - внутренняя ошибка сервера

### GET /{short_URL} - перенаправляет на оригинальный URL

Ожидаемое тело запроса:
```
http://localhost:8080/url-shortener-ozon/i9F9AvSymi
```

Возможные ответы:

* `301 Redirect` - короткая ссылка найдена, происходит редирект на оригинальный URL

* `404 Not Found` - короткая ссылка не найдена
```json
{
  "message": "resource not found",
  "code": "PS-00103"
}
```
* `500 Internal Server Error` - внутренняя ошибка сервера

---

### Алгоритм генерации сокращенных ссылок

1. Генерация случайной строки фиксированной длины (с использованием соли при коллизиях).
2. Проверка уникальности сгенерированного кода в хранилище.
3. При обнаружении коллизии - добавление соли и повторная генерация.

**Особенности:**
- Алгоритм гарантирует уникальность коротких ссылок;
- При повторной отправке одной и той же оригинальной ссылки всегда возвращается один и тот же короткий код;
- Поддержка двух режимов хранения: in-memory и PostgreSQL.

---

### Архитектура проекта

Проект построен с явным разделением на слои:

```
├── .github/             
│    └── workflows/
│        └── tests.yml              # GitHub Actions для запуска тестов
├── cmd/app/
│    └── main.go                    # Entrypoint - запуск приложения
├── internal/
│    ├── adapters/                  # Адаптеры для внешних систем
│    │   ├── controller/http/urlapi # Для работы со входящими и выходящими JSON структурами
│    │   └── repository/            # Для работы со структурами репозитория
│    ├── app/
│    │   └── app.go                 # Инициализация приложения
│    ├── apperror/                  # Кастомные ошибки приложения
│    ├── controller/                # Слой контроллеров (обработчики HTTP)
│    │   └── http/
│    │       └── v1/
│    │           └── urlshortener/  # Хендлеры для URL
│    ├── domain/                    # Сущности бизнес-логики
│    │   ├── entities/              # Модели данных
│    │   └── usecase/               # Слой бизнес-логики
│    │       ├── memory/            # Реализация сохранения во внутреннюю память
│    │       └── postgres/          # Реализация сохранения в базу данных
│    └── repository/url/
│        ├── memory/                # Инициализация и обращение во внутреннюю память
│        └── postgres/              # Инициализация и обращение в базы данных
├── pkg/
│    ├── config/                    # Конфигурация приложения
│    ├── connectors/                # Подключаемые модули
│    │   └── pgconnector/           # Подключение PostgreSQL
│    └── utils/                     # Вспомогательные функции
│        │── createlogger.go        # Инициализация логгера (zap)
│        │── generateresponse.go    # Генерация обращений
│        └── generateshorturl.go    # Генерация коротких ссылок
├── .config.yaml                    # Конфигурация приложения
├── .env.example                    # Шаблон переменных окружения
├── compose.yaml                    # Запуск через Docker
├── Dockerfile
├── go.mod                          # Зависимости проекта
└── urls.sql                        # Схема таблицы для PostgreSQL
```

**Ключевые архитектурные решения:**
- **Domain layer** - чистая бизнес-логика, независимая от внешних систем
- **Usecase layer** - оркестрация бизнес-процессов
- **Adapters layer** - адаптеры для HTTP и баз данных
- Поддержка двух профилей запуска через Docker Compose

---

### Установка и запуск

1. Клонировать репозиторий.
```bash
git clone https://github.com/RondaSMR/url-shortener-ozon.git
cd url-shortener-ozon
```

2. Создать файл .env из шаблона:
```bash
cp .env.example .env
# Отредактируйте .env, указав свои пароли
```

3. Запустить приложение в нужном режиме:

**Режим in-memory** (данные хранятся в оперативной памяти):
```bash
docker compose --profile memory up -d --build
```

**Режим PostgreSQL** (данные хранятся в базе данных):
```bash
docker compose --profile db up -d --build
```

Приложение будет доступно по адресу: `http://localhost:8080`

---

## Использование API

### Создание сокращенной ссылки

```bash
curl -i -X POST \
  -H "Content-Type: application/json" \
  -d '{"url": "https://ozon.ru"}' \
  "http://localhost:8080/url-shortener-ozon"

# HTTP/1.1 201 Created
# Content-Type: application/json; charset=utf-8
# Date: Fri, 06 Mar 2026 11:40:31 GMT
# Content-Length: 59
#
# {"success":true,"message":"","result":{"url":"i9F9AvSymi"}}%    
```

### Переход по короткой ссылке

**Через браузер:** просто введите в адресную строку:
```
http://localhost:8080/url-shortener-ozon/i9F9AvSymi
```

**Через curl** (без авто-редиректа, чтобы увидеть заголовки):
```bash
curl -i "http://localhost:8080/url-shortener-ozon/i9F9AvSymi"

# HTTP/1.1 301 Moved Permanently
# Content-Type: text/html; charset=utf-8
# Location: https://ozon.ru
# Date: Fri, 06 Mar 2026 11:45:19 GMT
# Content-Length: 50

# ...
```

**Примеры с ошибками**

Попытка создать ссылку с некорректным форматом:
```bash
curl -i -X POST \
  -H "Content-Type: application/json" \
  -d '{"url": "https/ozon.ru"}' \
  "http://localhost:8080/url-shortener-ozon"

# HTTP/1.1 400 Bad Request
# Content-Type: application/json; charset=utf-8
# Date: Fri, 06 Mar 2026 11:48:29 GMT
# Content-Length: 48
#
# {"message":"validation error","code":"PS-00102"}%   
```

Попытка получить несуществующую ссылку:
```bash
curl -i -X GET \
  -H "Content-Type: application/json" \
  -d '{"url": "i9F9FvSymi"}' \
  "http://localhost:8080/url-shortener-ozon"

# HTTP/1.1 404 Not Found
# Content-Type: text/plain
# Date: Fri, 06 Mar 2026 11:49:29 GMT
# Content-Length: 18
#
# 404 page not found%
```

---

## Тестирование

Проект покрыт тестами для обоих режимов хранения. Запуск тестов осуществляется через GitHub Actions при каждом push в main/master.

**Локальный запуск тестов:**

```bash
# Тесты для in-memory режима
go test -v ./internal/domain/usecase/memory/...

# Тесты для PostgreSQL режима
go test -v ./internal/domain/usecase/postgres/...
```

### Конфигурация

Основные параметры конфигурации задаются в `.config.yaml`:

```yaml
debug: true
serviceName: url-shortener-ozon
pgStorage:
  host: "postgres"
  port: 5432
http_server:
  address: "0.0.0.0:8080"
```

Чувствительные данные (пароли) вынесены в `.env` файл.

### Технологии

- **Go** 1.24+
- **PostgreSQL** 16 (для режима БД)
- **Docker** / **Docker Compose** - контейнеризация
- **GitHub Actions** - CI/CD
- **Gin** - HTTP фреймворк
- **pgx** - драйвер для PostgreSQL
- **zap** - логирование