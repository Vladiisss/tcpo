
# Практическая работа № 16

## Автор
Курков Владислав Николаевич
ПИМО-01-25
## Задание
Интеграционное тестирование API. Использование Docker для тестовой БД.

**Цели:**
- Освоить интеграционное тестирование REST API: проверка цепочки «маршрут → хендлер → сервис → репозиторий → реальная БД».
- Научиться поднимать изолированную тестовую среду БД в Docker.
- Освоить два подхода к инфраструктуре тестов:
    - **A.** Локальная среда через `docker-compose` (просто и наглядно).
    - **B.** Программный подъём контейнеров через `testcontainers-go` (изолированно и удобно для CI).
- Научиться инициализировать схему БД (миграции/auto-migrate), сидировать тестовые данные, очищать окружение.
- Внедрить интеграционные проверки CRUD-эндпоинтов (статусы, заголовки, JSON-ответы, эффекты в БД).

## Подготовка к запуску

### Клонирование репозитория
```bash
git clone https://github.com/MrFandore/Goland.git
cd Goland/Practica_16
```

### Структура проекта
```
Practica_16/
├── README.md
├── docker-compose.yml
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── db/
│   │   └── migrate.go
│   ├── httpapi/
│   │   └── handlers.go
│   ├── models/
│   │   └── note.go
│   ├── repo/
│   │   └── postgres.go
│   ├── service/
│   │   └── service.go
│   └── integration/
│       ├── notes_integration_test.go      # Тесты с docker-compose
│       └── notes_tc_integration_test.go   # Тесты с testcontainers-go
├── go.mod
└── go.sum
```

### Запуск приложения

#### Локальный запуск с Docker Compose
```bash
# Запуск PostgreSQL в Docker
docker-compose up -d

# Запуск приложения
go run ./cmd/api
```

## Примеры запросов к API

### 1. Создание заметки
```bash
curl -i -X POST http://178.72.139.210:8089/notes \
  -H 'Content-Type: application/json' \
  -d '{"title":"Hello","content":"World"}'
```
Результат:

![Создание заметки](Imagine/img.png)

### 2. Получение заметки по ID
```bash
curl -i http://178.72.139.210:8089/notes/1
```
Результат:

![Получение заметки](Imagine/img_1.png)

## Интеграционное тестирование

### Два подхода к тестовой инфраструктуре

#### Подход A: Локальная среда через docker-compose
Тесты используют предварительно поднятую через `docker-compose` БД. Подходит для локальной разработки.

```go
// notes_integration_test.go
func TestCreateNoteIntegration(t *testing.T) {
    // Подключение к БД, запущенной через docker-compose
    dsn := "postgres://test:test@localhost:5432/testdb?sslmode=disable"
    db, _ := sql.Open("postgres", dsn)
    
    repo := repo.NewPostgresRepo(db)
    service := service.NewService(repo)
    handler := httpapi.NewHandler(service)
    
    // Тестирование HTTP-эндпоинта
    req := httptest.NewRequest("POST", "/notes", body)
    w := httptest.NewRecorder()
    handler.CreateNote(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

#### Подход B: Программный подъём через testcontainers-go
Тесты сами поднимают и останавливают контейнер с БД. Идеально для CI/CD.

```go
// notes_tc_integration_test.go
func TestWithTestContainers(t *testing.T) {
    ctx := context.Background()
    
    // Запуск PostgreSQL в контейнере
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "test",
            "POSTGRES_PASSWORD": "test",
            "POSTGRES_DB":       "testdb",
        },
    }
    
    postgres, _ := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    defer postgres.Terminate(ctx)
    
    // Получение динамического порта
    port, _ := postgres.MappedPort(ctx, "5432")
    dsn := fmt.Sprintf("postgres://test:test@localhost:%s/testdb?sslmode=disable", port.Port())
    
    // Далее тесты аналогичны первому подходу
}
```

### Результаты запуска тестов

#### Запуск всех тестов
```bash
go test ./...
```
Результат:

![Запуск тестов](Imagine/img_2.png)

#### Запуск с подробным выводом
```bash
go test -v ./internal/integration/...
```
Результат:

![Подробный вывод тестов](Imagine/img_3.png)

## Выводы
В ходе практической работы были освоены подходы к интеграционному тестированию REST API с использованием реальной базы данных в Docker:

1. **Два подхода к инфраструктуре:**
    - `docker-compose` — простой и наглядный способ для локальной разработки, позволяет быстро поднять тестовое окружение.
    - `testcontainers-go` — более гибкий и изолированный подход, идеально подходящий для CI/CD, так как каждый тест получает чистое окружение.

2. **Проверка полной цепочки:** тесты охватывают взаимодействие всех слоёв приложения — от HTTP-запроса до изменений в базе данных, что даёт уверенность в корректной работе системы в целом.

3. **Управление тестовыми данными:** реализована инициализация схемы (миграции) и очистка данных между тестами, что обеспечивает изолированность и повторяемость тестов.

Интеграционные тесты являются критически важным компонентом качества для реальных проектов, так как они проверяют не только отдельные функции, но и корректность их взаимодействия с внешними зависимостями (БД, кэш и т.д.).

