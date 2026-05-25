# api-structure

REST API для управления организационной структурой: подразделения и сотрудники.

## Требования

- Docker
- Docker Compose
- Go 1.26.3+ для локального запуска и сборки

## Быстрый запуск

### 1. Поднять приложение

```bash
docker compose up -d --build
```

После успешного запуска сервис будет доступен по: http://localhost:8080

### Миграции
Миграции находятся в папке  migrations/  и представляют собой SQL-файлы goose.
При запуске через  docker compose up -d --build  миграции применяются автоматически отдельным сервисом `migrations` .

### PostgreSQL
Используются значения по умолчанию из  `docker-compose.yaml`:
- POSTGRES_USER=api_user 
- POSTGRES_PASSWORD=api_pass 
- POSTGRES_DB=api_structure 

---

## Технологии

- Go
- net/http
- GORM
- PostgreSQL
- goose
- Docker / Docker Compose

## Возможности

- Создание, изменение, удаление подразделений.
- Создание сотрудников внутри подразделений.
- Получение дерева подразделения до заданной глубины.
- Проверка ограничений:
  - нельзя создать сотрудника в несуществующем подразделении;
  - имя подразделения не пустое, длина 1..200;
  - имена подразделений уникальны внутри одного parent;
  - full_name и position не пустые, длина 1..200;
  - нельзя сделать подразделение родителем самого себя;
  - нельзя создать цикл в дереве;
  - удаление с cascade через БД/ORM;
  - при delete mode=reassign дети и сотрудники переносятся в целевой департамент.

---

## Структура проекта

```
api-structure
|- cmd
|   |-server
|   |   |-main.go
|- internal
|   |- models
|   |   |— models.go (модели GORM)
|   |- dto
|   |   |- dto.go (DTO для API)
|   |- http
|   |   |- handlers
|   |   |   |- handler.go
|   |   |   |- repository.go
|   |   |   |- service.go
|- migrations (QL-миграции goose)
|- docs
|   |- openapi.yaml
|- Dockerfile
|- docker-compose.yaml — запуск PostgreSQL, миграций и приложения.
```

## API

1) Создать подразделение
    POST /departments/
    body:
    ```
    {
        "name": "IT"
    }
    ```

2) Создать сотрудника
    POST /departments/{id}/employees/ 
    body:
    ```
    {
        "full_name": "Arseniy V",
        "position": "Backend Developer",
        "hired_at": "2026-05-05"
    }
    ```

3) Получить дерево подразделения
    GET /departments/{id}?depth=3&include_employees=true

4) Обновить подразделение
    PATCH /departments/{id} 
    body:
    ```
    {
        "name": "Platform",
        "parent_id": 2
    }
    ```

5) Удалить подразделение
    DELETE /departments/{id}?mode=cascade 
    или
    DELETE /departments/{id}?mode=reassign&reassign_to_department_id=2
