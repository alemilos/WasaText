package database

import (
	"database/sql"
	"fmt"
)

// Initialize the database tables if they don't exist. The initialized tables are:
// - users
func InitializeSchema(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			photo_path TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		// Add here other table cretions 
	}

	// all statements are executed in one transaction 
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	for _, s := range stmts {
		if _, err := tx.Exec(s); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("exec statement failed: %w\nstmt: %s", err, s)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}