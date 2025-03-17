package main

import (
	"fmt"
	"log"

	"github.com/BaikalMine/SongService/config"
	"github.com/BaikalMine/SongService/database"
	"github.com/BaikalMine/SongService/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения из .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Инициализация конфигурации
	cfg := config.LoadConfig()

	// Подключение к базе данных
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Запуск миграций (создание таблиц)
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Ошибка миграций: %v", err)
	}

	// Настройка маршрутов и запуск сервера
	router := routes.SetupRouter(db, cfg.ExternalAPIUrl)
	fmt.Printf("Сервер запущен на порту %s\n", cfg.Port)
	router.Run(":" + cfg.Port)
}
