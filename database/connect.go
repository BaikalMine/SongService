package database

import (
	"database/sql"
	"fmt"

	"github.com/BaikalMine/SongService/config"
	_ "github.com/lib/pq"
)

// Connect устанавливает подключение к PostgreSQL.
func Connect(cfg *config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
