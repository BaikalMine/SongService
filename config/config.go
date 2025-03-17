package config

import "os"

// Config хранит конфигурационные параметры приложения.
type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	ExternalAPIUrl string
	Port           string
}

// LoadConfig загружает конфигурацию из переменных окружения.
func LoadConfig() *Config {
	return &Config{
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		ExternalAPIUrl: os.Getenv("EXTERNAL_API_URL"),
		Port:           os.Getenv("PORT"),
	}
}
