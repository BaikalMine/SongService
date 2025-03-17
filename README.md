# SongService

SongService — это REST API для управления библиотекой песен, реализованное на Go. Проект предоставляет CRUD-операции для работы с песнями и позволяет получать текст песни с пагинацией. API интегрирован с Swagger для генерации документации, а работа с базой данных осуществляется с использованием транзакций для обеспечения целостности данных.

## Особенности

- **CRUD операции:** создание, получение, обновление и удаление песен.
- **Пагинация:** получение списка песен с поддержкой фильтрации и пагинации.
- **Обогащение данных:** при добавлении новой песни происходит запрос к внешнему API для получения дополнительной информации.
- **Swagger документация:** сгенерированная документация доступна через Swagger UI.
- **Транзакции:** все операции с базой данных (как для записи, так и для чтения) оборачиваются в транзакции.
- **Гибкая конфигурация:** параметры приложения загружаются из файла `.env`.

## Структура проекта
.
├── cmd
│   └── main.go            # Точка входа в приложение
├── config
│   └── config.go          # Загрузка конфигурационных параметров из .env
├── controllers
│   └── song_controller.go # HTTP-обработчики для работы с песнями
├── database
│   ├── connect.go         # Подключение к базе данных PostgreSQL
│   ├── migration.go       # Миграции (создание таблиц)
│   └── transaction.go     # Helper для выполнения транзакций
├── models
│   └── song.go            # Определение модели данных Song
├── routes
│   └── routes.go          # Определение маршрутов API
├── docs                   # Сгенерированная документация Swagger
├── go.mod                 # Файл зависимостей Go
└── .env                   # Файл конфигурации

## Требования

- Go 1.16 или выше
- PostgreSQL
- [swag](https://github.com/swaggo/swag) (опционально, для генерации документации)

## Установка и запуск

1. **Клонируйте репозиторий:**
   ```bash
   git clone https://github.com/BaikalMine/SongService.git
   cd SongService
2. **Установите зависимости:**
    go mod download
3.**Настройте переменные окружения:**