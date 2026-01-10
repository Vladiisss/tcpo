# Практическая работа № 9

## Автор
Курков Владислав Николаевич
ПИМО-01-25

## Задание

Реализация регистрации и входа пользователей. Хэширование паролей с bcrypt

Цели:
-	Научиться безопасно хранить пароли (bcrypt), валидировать вход и обрабатывать ошибки.
-	Реализовать эндпоинты POST /auth/register и POST /auth/login.
-	Закрепить работу с БД (PostgreSQL + GORM или database/sql) и валидацией ввода.
-	Подготовить основу для JWT-аутентификации в следующем ПЗ (№10). 

## Подготовка к запуску
Развёртывание БД
```bash
docker-compose -f docker-compose.dev.yml up -d
```
Установка зависимостей
```bash
make install
```
Запуск сервера
```bash
make run
```

### Конфигурация
.env
```golang
# Data Source Name (DSN) для подключения к базе данных
DB_DSN="host=localhost user=user password=pass dbname=pz9 port=5432 sslmode=disable"

# Имя базы данных PostgreSQL
POSTGRES_DB=pz9

# Имя пользователя PostgreSQL
POSTGRES_USER=user

# Пароль пользователя PostgreSQL
POSTGRES_PASSWORD=pass

# Порт, на котором запускается твое приложение (API)
APP_PORT=8080

# Внешний порт PostgreSQL для подключения снаружи контейнера
POSTGRES_EXTERNAL_PORT=5432

# Стоимость хэширования bcrypt (чем выше значение, тем медленнее и безопаснее)
BCRYPT_COST=12
```

## Фрагменты кода
Обработчик авторизации
```golang
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in loginReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json")
		return
	}
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	if in.Email == "" || in.Password == "" {
		writeErr(w, http.StatusBadRequest, "email_and_password_required")
		return
	}

	u, err := h.Users.ByEmail(context.Background(), in.Email)
	if err != nil {
		// не раскрываем, что именно не так
		writeErr(w, http.StatusUnauthorized, "invalid_credentials")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
		writeErr(w, http.StatusUnauthorized, "invalid_credentials")
		return
	}

	// В ПЗ10 здесь будет генерация JWT; пока просто ok
	writeJSON(w, http.StatusOK, authResp{
		Status: "ok",
		User:   map[string]any{"id": u.ID, "email": u.Email},
	})
}
```

Обработчик регистрации
```golang
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var in registerReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json")
		return
	}
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	if in.Email == "" || len(in.Password) < 8 {
		writeErr(w, http.StatusBadRequest, "email_required_and_password_min_8")
		return
	}

	// bcrypt hash
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), h.BcryptCost)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "hash_failed")
		return
	}

	u := core.User{Email: in.Email, PasswordHash: string(hash)}
	if err := h.Users.Create(r.Context(), &u); err != nil {

		if errors.Is(err, repo.ErrEmailTaken) {
			writeErr(w, http.StatusConflict, "email_taken")
			return
		}
		writeErr(w, http.StatusInternalServerError, "db_error")
		return
	}

	writeJSON(w, http.StatusCreated, authResp{
		Status: "ok",
		User:   map[string]any{"id": u.ID, "email": u.Email},
	})
}
```

## Результат работы

### 1. Результат автомиграции

![alt text](screenshots/image.png)

### 2. Регистрация
```bash
curl -i -X POST http://localhost:8080/auth/register  -H "Content-Type: application/json"  -d '{\"email\":\"user@example.com\",\"password\":\"Secret123!\"}'
```
Результат:

![alt text](screenshots/image-1.png)

### 3. Попытка регистрации (с существуюшим email)
```bash
curl -i -X POST http://localhost:8080/auth/register  -H "Content-Type: application/json"  -d '{\"email\":\"user@example.com\",\"password\":\"Secret123!\"}'
```
Результат:

![alt text](screenshots/image-2.png)


### 4. Авторизация
```bash
curl.exe -i -X POST http://localhost:8080/auth/login  -H "Content-Type: application/json"  -d '{\"email\":\"user@example.com\",\"password\":\"Secret123!\"}'
```
Результат:

![alt text](screenshots/image-3.png)

### 5. Авторизация (с неверными данными)
```bash
curl -Method DELETE http://localhost:8080/api/v1/notes/6904e846613fbf31ddac61e5
```
Результат:

![alt text](screenshots/image-4.png)


## Выводы 

Практическая работа позволила освоить ключевой аспект безопасности веб-приложений — безопасное хранение паролей с использованием bcrypt для хэширования и сравнения, что стало базой для понимания принципов защиты пользовательских данных от утечек. Реализация эндпоинтов регистрации (`POST /auth/register`) и входа (`POST /auth/login`) дала опыт полной цепочки обработки пользовательских данных: от валидации входных параметров и проверки уникальности (email/username) до создания записей в PostgreSQL через GORM/database/sql с автоматической обработкой ошибок дублирования и некорректного ввода.

Полученная структура аутентификации — с проверкой существования пользователя, сравнением хэшей паролей и возвратом соответствующих HTTP-кодов (201 Created, 200 OK, 400 Bad Request, 401 Unauthorized) — соответствует промышленным стандартам базовой авторизации и подготавливает почву для следующего этапа с JWT-токенами. Работа с валидацией ввода (структуры с тегами `validate`) и обработкой типичных сценариев (существующий пользователь, неверный пароль, пустые поля) закрепила навыки построения надёжного слоя аутентификации, готового к интеграции с middleware и сессионным управлением в дальнейшем развитии проекта.
