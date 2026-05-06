# URL Shortener

Сервис для сокращения URL-адресов с использованием Go, PostgreSQL и Chi роутера.

## Описание

Этот проект представляет собой REST API для создания, получения, обновления и удаления сокращенных URL. Сервис поддерживает:

- Создание коротких ссылок
- Перенаправление на оригинальный URL
- Обновление существующих ссылок
- Удаление ссылок
- Аутентификацию для защищенных операций

## Функции

- **Сокращение URL**: Создание коротких алиасов для длинных URL
- **Перенаправление**: Автоматическое перенаправление по короткому алиасу
- **Обновление**: Изменение URL для существующего алиаса
- **Удаление**: Удаление коротких ссылок
- **Аутентификация**: Защита операций обновления и удаления

## Требования

- Go 1.21+
- PostgreSQL 12+
- Git

## Установка

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/your-username/url-shortener.git
   cd url-shortener
   ```

2. Установите зависимости:
   ```bash
   go mod tidy
   ```

3. Настройте базу данных PostgreSQL и переменные окружения.

## Настройка

### 1. База данных

Создайте базу данных PostgreSQL. Таблица `url` будет создана автоматически при запуске приложения.

### 2. Переменные окружения

Создайте файл `.env` в корне проекта:

```env
CONFIG_PATH=./config/local.yaml
DATABASE_URL=postgres://username:password@localhost:5432/url_shortener?sslmode=disable
user=admin
HTTP_SERVER_PASSWORD=your_secure_password
```

### 3. Конфигурационный файл

Создайте файл `config/local.yaml`:

```yaml
env: "local"
http_server:
  address: "localhost:8082"
  timeout: 4s
  idle_timeout: 60s
```

## Запуск

Запустите сервер:

```bash
go run cmd/url-shortener/main.go
```

Сервер будет доступен по адресу `http://localhost:8082`.

## Использование API

### Создание короткой ссылки

**POST /url**

Создает новую короткую ссылку. Если алиас не указан, генерируется автоматически.

**Запрос:**
```json
{
  "url": "https://example.com/very/long/url",
  "alias": "my-link"
}
```

**Ответ:**
```json
{
  "status": "OK",
  "alias": "my-link"
}
```

### Получение оригинального URL

**GET /{alias}**

Перенаправляет на оригинальный URL.

**Пример:**
```
GET http://localhost:8082/my-link
```

### Обновление ссылки

**PUT /url**

Обновляет URL для существующего алиаса. Требует аутентификации.

**Запрос:**
```json
{
  "alias": "my-link",
  "new_url": "https://new-example.com"
}
```

**Ответ:**
```json
{
  "status": "OK"
}
```

### Удаление ссылки

**DELETE /{alias}**

Удаляет короткую ссылку. Требует аутентификации.

**Пример:**
```
DELETE http://localhost:8082/my-link
```

**Ответ:**
```json
{
  "status": "OK"
}
```

## Аутентификация

Операции `PUT /url` и `DELETE /{alias}` защищены Basic Auth. Используйте учетные данные из переменных окружения:

- **Username**: значение переменной `user`
- **Password**: значение переменной `HTTP_SERVER_PASSWORD`

## Примеры использования

### Создание ссылки через curl

```bash
curl -X POST http://localhost:8082/url \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com", "alias": "github"}'
```

### Обновление ссылки через curl

```bash
curl -X PUT http://localhost:8082/url \
  -H "Content-Type: application/json" \
  -u admin:your_password \
  -d '{"alias": "github", "new_url": "https://github.com/features"}'
```

### Проверка перенаправления

```bash
curl -I http://localhost:8082/github
```

## Структура проекта

```
url-shortener/
├── cmd/url-shortener/     # Точка входа приложения
├── config/                # Конфигурационные файлы
├── internal/
│   ├── config/           # Загрузка конфигурации
│   ├── http-server/      # HTTP сервер и обработчики
│   ├── lib/              # Вспомогательные библиотеки
│   └── storage/          # Хранилище данных
├── tests/                 # Интеграционные тесты
└── go.mod                 # Go модули
```

## Разработка

### Запуск тестов

```bash
go test ./...
```

