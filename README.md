# Сокращатель ссылок

---

Сервер принимает два запроса:

**GET /** - получение оригинальной ссылки по сокращенной

Ожидаемое тело запроса:
```json
{
  "url": "i9F9AvSymi"
}
```

Возможные ответы:
    
* `200 OK` - короткая ссылка найдена, оригинальная ссылка возвращена
```json
{
  "success": true,
  "message": "",
  "result": {
    "url": "https://ozon.ru"
  }
}
```
* `404 Not Found` - короткая ссылка не найдена
```json
{
    "message": "resource not found",
    "code": "PS-00103"
}
```
* `500 Internal Server Error` - внутренняя ошибка сервера


**POST /** - создание сокращенной ссылки

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

---

**Алгоритм генерации сокращенных ссылок**

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

### Использование API

**Создание сокращенной ссылки**

```bash
curl -i -X POST \                       
  -u myuser:mypass \
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

**Получение оригинальной ссылки**

```bash
curl -i -X GET \
  -u myuser:mypass \
  -H "Content-Type: application/json" \
  -d '{"url": "i9F9AvSymi"}' \
  "http://localhost:8080/url-shortener-ozon"

# HTTP/1.1 200 OK
# Content-Type: application/json; charset=utf-8
# Date: Fri, 06 Mar 2026 11:46:57 GMT
# Content-Length: 64
#
# {"success":true,"message":"","result":{"url":"https://ozon.ru"}}%    
```

**Примеры с ошибками**

Попытка создать ссылку с некорректным форматом:
```bash
curl -i -X POST \                      
  -u myuser:mypass \                         
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
  -u myuser:mypass \
  -H "Content-Type: application/json" \
  -d '{"url": "i9F9FvSymi"}' \  
  "http://localhost:8080/url-shortener-ozon"

# HTTP/1.1 404 Not Found
# Content-Type: application/json; charset=utf-8
# Date: Fri, 06 Mar 2026 11:49:29 GMT
# Content-Length: 50
#
# {"message":"resource not found","code":"PS-00103"}%  
```

### Тестирование

Проект покрыт тестами для обоих режимов хранения. Запуск тестов осуществляется через GitHub Actions при каждом push в main/master.

**Локальный запуск тестов:**

```bash
# Тесты для in-memory режима
STORAGE_MODE=memory go test ./...

# Тесты для PostgreSQL режима
STORAGE_MODE=db go test ./...
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