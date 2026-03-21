# Department & Employee Management API

REST API для управления иерархической структурой подразделений и сотрудников.

---

## 🚀 Быстрый старт (Docker)

### 1️⃣ Требования
- Docker & Docker Compose

### 2️⃣ Подготовка
```bash
# Клонируйте репозиторий
git clone <your-repo>
cd testtt

# (Опционально) Если нужны другие переменные окружения:
# Отредактируйте docker-compose.yml переменные окружения
```

### 3️⃣ Запуск
```bash
docker-compose up --build
```

**Ожидаемый результат:**
- PostgreSQL запустится на `localhost:5432`
- Миграции применятся автоматически
- API будет доступен на `http://localhost:8080`

### 4️⃣ Проверка
```bash
# Создать подразделение
curl -X POST http://localhost:8080/departments/ \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Department"}'

# Получить подразделение
curl "http://localhost:8080/departments/1?depth=1&include_employees=true"
```

---

## 📡 API Endpoints

| Метод | Endpoint | Описание |
|-------|----------|---------|
| POST | `/departments/` | Создать подразделение |
| POST | `/departments/{id}/employees/` | Добавить сотрудника |
| GET | `/departments/{id}` | Получить подразделение с деревом |
| PATCH | `/departments/{id}` | Обновить подразделение |
| DELETE | `/departments/{id}` | Удалить подразделение |

**Query параметры для GET:**
- `depth` (int, default=1, max=5) - глубина дерева
- `include_employees` (bool, default=true) - показывать сотрудников

**Query параметры для DELETE:**
- `mode` (cascade|reassign) - режим удаления
- `reassign_to_department_id` (int) - требуется если mode=reassign

---

## 💡 Примеры запросов

```bash
# 1. Создать подразделение
curl -X POST http://localhost:8080/departments/ \
  -H "Content-Type: application/json" \
  -d '{"name":"Engineering"}'

# 2. Добавить сотрудника (в department_id=1)
curl -X POST http://localhost:8080/departments/1/employees/ \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "position": "Developer",
    "hired_at": "2024-01-01"
  }'

# 3. Получить дерево
curl "http://localhost:8080/departments/1?depth=3&include_employees=true"

# 4. Обновить подразделение
curl -X PATCH http://localhost:8080/departments/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Tech Division"}'

# 5. Удалить с каскадом
curl -X DELETE "http://localhost:8080/departments/1?mode=cascade"

# 6. Удалить и перевести сотрудников
curl -X DELETE "http://localhost:8080/departments/1?mode=reassign&reassign_to_department_id=2"
```

---

## 🧪 Полный набор тестов

Все примеры запросов (включая edge cases и ошибки) находятся в **[API_TESTS.md](API_TESTS.md)**

---

## 🛢️ Переменные окружения

**В docker-compose.yml:**
```yaml
environment:
  SERVER_PORT=8080
  POSTGRES_HOST=postgres
  POSTGRES_PORT=5432
  POSTGRES_USER=gigi
  POSTGRES_PASSWORD=gogi
  POSTGRES_DB=go
  POSTGRES_CONNS=10
```

**Для локальной разработки без Docker** создайте `.env`:
```env
SERVER_PORT=8080
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=gigi
POSTGRES_PASSWORD=gogi
POSTGRES_DB=go
POSTGRES_CONNS=10
```

---

## 🗄️ Архитектура БД

```sql
-- Подразделения (иерархия через self-reference)
department (
  id: int PRIMARY KEY,
  name: varchar(200) NOT NULL,
  parent_id: int NULL FK→department(id) ON DELETE CASCADE,
  created_at: timestamp DEFAULT NOW(),
  UNIQUE (name, parent_id)  -- уникально в пределах parent
)

-- Сотрудники
employee (
  id: int PRIMARY KEY,
  department_id: int NOT NULL FK→department(id) ON DELETE CASCADE,
  full_name: varchar(200) NOT NULL,
  position: varchar(200) NOT NULL,
  hired_at: date NULL,
  created_at: timestamp DEFAULT NOW()
)
```

---

## 📋 Функциональность

✅ Иерархические подразделения (дерево)  
✅ Управление сотрудниками  
✅ Каскадное удаление (ON DELETE CASCADE)  
✅ Переназначение при удалении (reassign mode)  
✅ Валидация данных (длина, пустые значения)  
✅ Проверка циклов в иерархии  
✅ Уникальность имен в пределах parent  
✅ FK проверка (404 для несуществующих)  
✅ Структурированное логирование  
✅ Unit тесты  

---

## ✅ Валидация и ограничения

| Параметр | Правило | Статус код |
|----------|---------|-----------|
| name (department) | 1-200 символов, не пустой, триммированный | 400 |
| full_name (employee) | 1-200 символов, не пустой | 400 |
| position | 1-200 символов, не пустой | 400 |
| Дублирующееся имя | Не допускается в одном parent | 409 |
| Циклы в иерархии | Не допускаются | 409 |
| Самоссылка | parent_id ≠ id | 409 |
| FK на department | Должно существовать | 404 |
| GET несуществующего | Возвращает 404 | 404 |

---

## 🏗️ Структура проекта

```
.
├── cmd/main.go                    # Точка входа
├── internal/
│   ├── handlers/                  # HTTP обработчики
│   │   ├── handlers.go
│   │   └── handlers_test.go
│   ├── middleware/                # Маршрутизация
│   ├── models/                    # Структуры данных
│   └── storage/postgre/           # CRUD + бизнес-логика
├── pkg/
│   ├── config/                    # Конфигурация из .env
│   └── logger/                    # Логирование slog
├── migrations/
│   ├── 000001_init.up.sql
│   └── 000001_init.down.sql
├── docker-compose.yml
├── Dockerfile
├── .env                           # (создается при запуске)
└── API_TESTS.md                   # Примеры всех запросов
```

---

## 🐛 Отладка

```bash
# Посмотреть логи
docker-compose logs -f golang

# Посмотреть логи PostgreSQL
docker-compose logs -f postgres

# Перезапустить сервис
docker-compose restart golang

# Полный рестарт (очистить data)
docker-compose down -v
docker-compose up --build

# Удалить контейнеры полностью
docker-compose down
```

---

## 🧩 Технический стек

- **Language:** Go 1.25.7
- **HTTP:** net/http (stdlib)
- **ORM:** GORM
- **Database:** PostgreSQL
- **Migrations:** golang-migrate (docker image)
- **Logging:** slog + tint
- **Docker:** Docker & Docker Compose
- **Tests:** standard testing + httptest

---

## 📄 Лицензия

MIT
