package database

import (
	"database/sql"
)

// WithTransaction выполняет функцию fn в рамках транзакции.
// Если fn возвращает ошибку, транзакция откатывается, иначе — коммитится.
func WithTransaction(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
