
# Практическая работа № 11

## Автор
Курков Владислав Николаевич
ПИМО-01-25

## Задание
Проектирование REST API (CRUD для заметок). Разработка структуры.

**Цели:**
- Освоить принципы проектирования REST API.
- Спроектировать и реализовать CRUD-интерфейс (Create, Read, Update, Delete) для сущности «Заметка».
- Подготовить основу для интеграции с базой данных и JWT-аутентификацией.

## Подготовка к запуску

### Инициализация проекта и установка зависимостей
```bash
mkdir notes-api
cd notes-api
go mod init example.com/notes-api
go get github.com/go-chi/chi/v5
```

### Структура проекта
```bash
Prak_11/
├── cmd/
│   └── api/
│       └── main.go                 # Точка входа в приложение. Инициализирует репозитории, обработчики и запускает HTTP-сервер.
├── internal/                       # Приватный код, который не должен импортироваться внешними проектами.
│   ├── api/
│   │   └── openapi.yaml            # Спецификация REST API в формате OpenAPI (Swagger).
│   ├── core/                       # Слой бизнес-сущностей и интерфейсов.
│   │   ├── note.go                 # Определение структуры (модели) Note.
│   │   └── service/                # Слой бизнес-логики.
│   │       └── note_service.go     # Реализация сервиса для работы с заметками.
│   ├── http/                       # Слой HTTP/API.
│   │   ├── handlers/               # Обработчики HTTP-запросов.
│   │   │   └── notes.go            # Реализация CRUD-обработчиков для ресурса /notes.
│   │   └── router.go               # Настройка маршрутизации (роутер).
│   └── repo/                       # Слой данных (Repository).
│       └── note_mem.go             # Реализация репозитория заметок в оперативной памяти (In-Memory).
└── go.mod                          # Файл управления зависимостями.
```

### Код приложения

#### Модель данных (internal/core/note.go)
```go
package core

import "time"

type Note struct {
    ID        int64
    Title     string
    Content   string
    CreatedAt time.Time
    UpdatedAt *time.Time
}
```

#### In-memory репозиторий (internal/repo/note_mem.go)
```go
package repo

import (
    "sync"
    "example.com/notes-api/internal/core"
)

type NoteRepoMem struct {
    mu    sync.Mutex
    notes map[int64]*core.Note
    next  int64
}

func NewNoteRepoMem() *NoteRepoMem {
    return &NoteRepoMem{notes: make(map[int64]*core.Note)}
}

func (r *NoteRepoMem) Create(n core.Note) (int64, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.next++
    n.ID = r.next
    r.notes[n.ID] = &n
    return n.ID, nil
}
```

#### HTTP-обработчик (internal/http/handlers/notes.go)
```go
package handlers

import (
    "encoding/json"
    "net/http"
    "example.com/notes-api/internal/core"
    "example.com/notes-api/internal/repo"
)

type Handler struct {
    Repo *repo.NoteRepoMem
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
    var n core.Note
    if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    id, _ := h.Repo.Create(n)
    n.ID = id
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(n)
}
```

#### Маршрутизация (internal/http/router.go)
```go
package httpx

import (
    "github.com/go-chi/chi/v5"
    "example.com/notes-api/internal/http/handlers"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
    r := chi.NewRouter()
    r.Post("/api/v1/notes", h.CreateNote)
    return r
}
```

#### Точка входа (cmd/api/main.go)
```go
package main

import (
    "log"
    "net/http"
    "example.com/notes-api/internal/http"
    "example.com/notes-api/internal/http/handlers"
    "example.com/notes-api/internal/repo"
)

func main() {
    repo := repo.NewNoteRepoMem()
    h := &handlers.Handler{Repo: repo}
    r := httpx.NewRouter(h)

    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
```

### Запуск сервера
```bash
go run ./cmd/api
```
Результат:

![Запуск приложения](foto/start_app.png)

### Пример запроса

#### Создание заметки (POST /api/v1/notes)
```bash
curl -X POST http://localhost:8080/api/v1/notes \
-H "Content-Type: application/json" \
-d '{"title":"Первая заметка", "content":"Это тест"}'
```
Результат:

![Создание заметки](foto/create_note.png)

## Контрольные вопросы

### 1. Что означает аббревиатура REST и в чём её суть?

**REST** расшифровывается как **RE**presentational **S**tate **T**ransfer (Передача репрезентативного состояния).

**Суть:** Это архитектурный стиль для создания распределённых систем, таких как веб-сервисы. Взаимодействие между клиентом и сервером происходит вокруг **ресурсов** (например, notes, users). Клиент взаимодействует с ресурсом, используя стандартные методы HTTP, и получает его **представление** (обычно в формате JSON или XML), после чего переходит в новое состояние.

---

### 2. Как связаны CRUD-операции и методы HTTP?

CRUD-операции напрямую сопоставляются с HTTP-методами:
- **Create** (Создание) → `POST`
- **Read** (Чтение) → `GET`
- **Update** (Обновление) → `PUT` / `PATCH`
- **Delete** (Удаление) → `DELETE`

---

### 3. Для чего нужна слоистая архитектура (handler → service → repository)?

Слоистая архитектура нужна для **разделения ответственности**. Это даёт несколько преимуществ:
- **Тестирование:** Каждый слой можно тестировать независимо.
- **Гибкость:** Можно заменить реализацию одного слоя (например, перейти с in-memory репозитория на PostgreSQL), не меняя другие слои.
- **Поддерживаемость:** Код становится более организованным и понятным, что упрощает его поддержку и развитие.

---

### 4. Что означает принцип «stateless» в REST API?

**Stateless** (без состояния) означает, что сервер **не хранит информацию о сессии** клиента между запросами. Каждый запрос от клиента к серверу должен содержать **всю необходимую информацию** для его полной обработки (например, данные аутентификации в заголовке `Authorization`).

---

### 5. Почему важно использовать стандартные коды ответов HTTP?

Использование стандартных кодов (2xx, 4xx, 5xx) обеспечивает:
- **Единообразие:** Все разработчики понимают их одинаково.
- **Предсказуемость:** Клиент может однозначно определить результат запроса: успех (2xx), ошибка на стороне клиента (4xx) или ошибка сервера (5xx).

---

### 6. Как можно добавить аутентификацию в REST API?

Наиболее популярный способ для REST API — использование **Bearer Токенов**, чаще всего **JWT** (JSON Web Tokens). После успешного логина клиент получает токен и передаёт его в каждом запросе в заголовке:
```
Authorization: Bearer <token>
```

---

### 7. В чём преимущество версионирования API (например, `/api/v1/`)?

Версионирование позволяет **развивать API** и вносить несовместимые изменения (breaking changes) без нарушения работы **старых клиентов**. Старые клиенты продолжают использовать предыдущую версию (`/api/v1/`), в то время как новые могут использовать функционал новой версии (`/api/v2/`).

## Выводы
В ходе практической работы был спроектирован и реализован базовый CRUD-интерфейс для сущности «Заметка» в соответствии с принципами REST. Была применена слоистая архитектура (handler, repo), что позволило разделить ответственность между обработкой HTTP-запросов, бизнес-логикой и хранением данных.

На практике были освоены принципы проектирования REST API, включая использование корректных HTTP-методов, статус-кодов и структуры URL. Полученная кодовая база заложила основу для дальнейшего развития проекта: интеграции с реальной базой данных и добавления JWT-аутентификации.