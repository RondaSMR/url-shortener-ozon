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

Ожидаемый URL:
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
│        └── tests.yml               # GitHub Actions для запуска тестов
├── cmd/app/
│    └── main.go                     # Entrypoint - запуск приложения
├── internal/
│    ├── app/
│    │   └── app.go                  # Инициализация приложения
│    ├── apperror/                   # Кастомные ошибки приложения
│    ├── controller/                 # Слой контроллеров (обработчики HTTP)
│    │   └── http/
│    │       ├── entities/           # DTO модели данных для HTTP слоя
│    │       └── v1/
│    │           └── urlshortener/   # Хендлеры для URL
│    ├── domain/                     # Сущности бизнес-логики
│    │   ├── entities/               # Модели данных
│    │   └── usecase/                # Слой бизнес-логики
│    │       ├── url_usecase.go      # Реализация бизнес-логики
│    │       └── url_usecase_test.go # Тестирование бизнес-логики
│    └── repository/url/
│        ├── entities/               # SQL модели данных
│        └── url/
│            ├── memory/
│            │   ├── init.go         # Инициализация in-memory репозитория
│            │   └── repository.go   # Обращение в репозитрий Postgres
│            └── postgres/
│                ├── init.go         # Инициализация Postgres репозитория
│                └── repository.go   # Обращение в in-memory репозитрий
├── pkg/
│    ├── config/                     # Конфигурация приложения
│    ├── connectors/                 # Подключаемые модули
│    │   └── pgconnector/            # Подключение PostgreSQL
│    └── utils/                      # Вспомогательные функции
│        │── create_logger.go        # Инициализация логгера (zap)
│        │── generate_response.go    # Генерация обращений
│        └── generate_short_url.go   # Генерация коротких ссылок
├── .config.yaml                     # Конфигурация приложения
├── .env.example                     # Шаблон переменных окружения
├── compose.yaml                     # Запуск через Docker
├── Dockerfile
├── go.mod                           # Зависимости проекта
└── urls.sql                         # Схема таблицы для PostgreSQL
```

**Ключевые архитектурные решения:**
- **Delivery layer** - обработка HTTP-запросов (контроллеры, middleware);
- **Domain layer** - чистая бизнес-логика, независимая от внешних систем;
- **Usecase layer** - оркестрация бизнес-процессов;
- **Repository layer** - работа с PostgreSQL и in-memory хранилищем (репозитории);
- Поддержка двух профилей запуска через Docker Compose.

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
# Отредактируйте .env, указав свои переменные окружения
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

Попытка перейти по несуществующей ссылке:
```bash
curl -i "http://localhost:8080/url-shortener-ozon/i9F9AvSymi"

# HTTP/1.1 404 Not Found
# Content-Type: text/plain
# Content-Type: application/json; charset=utf-8
# Content-Length: 50
#
# {"message":"resource not found","code":"PS-00103"}%
```

---

## Тестирование

Проект покрыт тестами для обоих режимов хранения. Запуск тестов осуществляется через GitHub Actions при каждом push в main/master.

**Локальный запуск тестов:**

```bash
# Тесты для in-memory режима
go test -v ./internal/domain/usecase/... -run WithMemory

# Тесты для PostgreSQL режима
go test -v ./internal/domain/usecase/... -run WithPostgres
```

### Конфигурация

Отладочные параметры конфигурации задаются в `.config.yaml`:

```yaml
debug: true
serviceName: url-shortener-ozon
```

Основные данные вынесены в `.env` файл.

```dotenv
POSTGRES_HOST=your_host
POSTGRES_PORT=your_port
POSTGRES_USER=your_db_user
POSTGRES_PASSWORD=your_db_password
POSTGRES_DB=urls

HTTP_SERVER_ADDRESS=your_http_server_address
```

### Технологии

- **Go** 1.24+
- **PostgreSQL** 16 (для режима БД)
- **Docker** / **Docker Compose** - контейнеризация
- **GitHub Actions** - CI/CD
- **Gin** - HTTP фреймворк
- **pgx** - драйвер для PostgreSQL
- **zap** - логирование
- 