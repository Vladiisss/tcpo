# Практическая работа № 6

## Автор
Курков Владислав Николаевич
ПИМО-01-25

## Задание
Использование ORM (GORM). Модели, миграции и связи между таблицами.

Цели:
- Понять, что такое ORM и чем удобен GORM.
- Научиться описывать модели Go-структурами и автоматически создавать таблицы (миграции через AutoMigrate).
- Освоить базовые связи: 1:N и M:N + выборки с Preload.
- Написать короткий REST (2–3 ручки) для проверки результата.

## Подготовка к запуску

### Установка
```bash
make install
```

### Отладка
```bash
make run
```

### Билд
```bash
make build
```

```bash
.\server
```

## Конфигурация
Переменные окружения:
- DB_DSN - строка для подключения к БД

## Окружение
- Сервер: go 1.23.0
- БД: PostgreSQL 16.10

## Примеры запросов

### 1. Примеры запросов
```bash
curl "http://localhost:8080/health"
```
Результат:

![alt text](docs/image.png)

### 2. Создание пользователя
```bash
curl -Uri "http://localhost:8080/users" `
  -Method POST `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "name": "Max",
    "email": "max@example.com"
  }'

```
Результат:

![alt text](docs/image-1.png)

### 3. Создание заметки
```bash
curl -Uri "http://localhost:8080/notes" `
  -Method POST `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{
    "title": "First note",
    "content": "Example content text",
    "userId": 1,
    "tags": ["go", "gorm", "practice"]
  }'

```
Результат:

![alt text](docs/image-2.png)

### 4. Получение заметки по ID
```bash
curl -Uri "http://localhost:8080/notes/1" -Method GET
```
Результат:

![alt text](docs/image-3.png)

## Контрольные вопросы

### 1. ORM и database/sql

ORM (Object-Relational Mapping) — библиотека для преобразования объектов Go в SQL-запросы и обратно, упрощая работу с БД.[1]
Нужна помимо database/sql для автоматизации boilerplate-кода: миграций, связей, валидации. **Плюсы**: сокращает SQL-код, поддерживает типобезопасные запросы, авто-миграции (AutoMigrate). **Минусы**: скрывает SQL (сложно дебажить), overhead на производительность, проблемы с сложными запросами.[2][3][4]

### 2. Связи 1:N и M:N в GORM
1:N (один ко многим):
В модели User поле Notes []Note создаёт связь один пользователь → много заметок.

M:N (многие ко многим):
В модели Note поле Tags []Tag 'gorm:"many2many:note_tags;"' и в модели Tag поле Notes []Note'gorm:"many2many:note_tags;"' создают связь много заметок ↔ много тегов через таблицу note_tags.

Пример моделей:

```golang
type User struct {
    ID    uint
    Name  string
    Email string
    Notes []Note
}

type Note struct {
    ID      uint
    Title   string
    Content string
    UserID  uint
    User    User
    Tags    []Tag `gorm:"many2many:note_tags;"`
}

type Tag struct {
    ID    uint
    Name  string
    Notes []Note `gorm:"many2many:note_tags;"`
}
```

### 3. AutoMigrate

AutoMigrate сканирует структуры Go и создает/обновляет таблицы, индексы, FK, constraints.[3]
Недостаточно для: сложных миграций (drop columns/tables), data-мigration, production (нет версионирования, риск потери данных).[4][2]

### 4. Preload vs Find/First

Preload загружает связанные данные (associations) отдельными запросами, избегая N+1 проблемы (bulk IN-запросы).[5][6]
Find/First берут только основную модель, связи — lazy или nil. Применять Preload при выборке с relations: `db.Preload("User.Tags").Find(&notes)`.[7]

### 5. Обработка unique index в GORM

При нарушении unique constraint PostgreSQL возвращает pgconn.PgError с Code="23505" (unique_violation) и ConstraintName.[8][9]
В GORM: проверьте `errors.As(err, &pgError)` и сравните `pgError.Code == pgconn.PgErrorCodeUniqueViolation` или имя constraint.[8]


## Выводы
Освоено создание моделей через структуры Go с тегами gorm, автоматические миграции с AutoMigrate и связи 1:N (foreignKey) и M:N (many2many). Реализованы REST-эндпоинты для создания/получения данных с Preload для связей, что устраняет N+1 проблему.