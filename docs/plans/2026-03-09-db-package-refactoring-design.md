# Рефакторинг пакета internal/db

## Scope

Только пакет `internal/db`. Остальные слои (portal, rest, domain) не затрагиваются.

## Модели (models.go)

| Текущее          | Новое      | Изменения                         |
|------------------|------------|-----------------------------------|
| News             | News       | Без изменений                     |
| NewsListItemDAO  | ShortNews  | Переименование                    |
| TagDAO           | Tag        | Переименование                    |
| —                | Category   | Новая: ID, Name, SortOrder, StatusID |

## Сигнатуры репозитория

Возвращаемые типы — db-модели вместо domain:

- `ListNews(...)` → `([]ShortNews, error)`
- `GetNews(...)` → `(*News, error)`
- `GetCategories(...)` → `([]Category, error)`
- `GetTags(...)` → `([]Tag, error)`
- `getTagsByIDs(...)` → `([]Tag, error)`
- `getTagsMapByIDs(...)` → `(map[int32]Tag, error)`
- `CountNews(...)` → без изменений (`int`)

## Изменения в логике

- В ListNews и GetNews теги возвращаются как []Tag, без конвертации в domain
- ShortNews и News хранят теги как []Tag
- GetCategories сканирует все 4 поля (id, name, sortOrder, statusId)

## Что НЕ меняется

- Входные параметры (domain.ListNewsReq и т.д.)
- SQL-запросы (кроме GetCategories)
- Структура файлов пакета
