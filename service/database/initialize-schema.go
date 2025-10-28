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
		`CREATE TABLE IF NOT EXISTS conversations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL CHECK(type IN ('private','group')),
			photo_path TEXT,	
			name TEXT, 
			created_by INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
			-- constraints for type 'group'
			CHECK (type != 'group' OR name IS NOT NULL)
		);`,
		`CREATE TABLE IF NOT EXISTS conversation_members (
			conversation_id INTEGER NOT NULL,
    		user_id INTEGER NOT NULL,
    		role TEXT, -- "admin" or "member", valid only for group conversations
    		joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    		PRIMARY KEY (conversation_id, user_id),
    		FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE	
		);`,
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
    		type TEXT NOT NULL CHECK(type IN ('text', 'image')),
    		is_forwarded INTEGER NOT NULL DEFAULT 0 CHECK(is_forwarded IN (0,1)),
    		conversation_id INTEGER NOT NULL,
    		author_id INTEGER NOT NULL,
    		content TEXT NOT NULL,
			secondary_content TEXT NOT NULL,
			reply_to INTEGER, 
    		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    		FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    		FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS message_read(
			message_id INTEGER NOT NULL,
			member_id INTEGER NOT NULL,
			read_at DATETIME DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (message_id, member_id),
			FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
			FOREIGN KEY (member_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS comments (
			message_id INTEGER NOT NULL,
			author_id INTEGER NOT NULL,
			emoji TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (message_id, author_id),
			FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
			FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
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
