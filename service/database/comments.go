package database

import "database/sql"

type Comment struct {
	MessageID int64  `json:"message_id"`
	AuthorID  int64  `json:"author_id"`
	Emoji     string `json:"emoji"`
	CreatedAt string `json:"created_at"`
}

func (db *appdbimpl) GetCommentByUser(messageID, authorID int64) (*Comment, error) {
	row := db.c.QueryRow("SELECT message_id, author_id, emoji, created_at FROM comments WHERE message_id = ? AND author_id = ?", messageID, authorID)
	var c Comment
	err := row.Scan(&c.MessageID, &c.AuthorID, &c.Emoji, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (db *appdbimpl) GetCommentsByMessage(messageID int64) ([]Comment, error) {
	rows, err := db.c.Query(`
		SELECT message_id, author_id, emoji, created_at
		FROM comments
		WHERE message_id = ?
		ORDER BY created_at ASC;
	`, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.MessageID, &c.AuthorID, &c.Emoji, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (db *appdbimpl) AddOrUpdateComment(messageID, authorID int64, emoji string) error {
	_, err := db.c.Exec(`
		INSERT INTO comments (message_id, author_id, emoji)
		VALUES (?, ?, ?)
		ON CONFLICT(message_id, author_id)
		DO UPDATE SET emoji = excluded.emoji, created_at = CURRENT_TIMESTAMP;
	`, messageID, authorID, emoji)
	return err
}

func (db *appdbimpl) DeleteComment(messageID, authorID int64) error {
	_, err := db.c.Exec("DELETE FROM comments WHERE message_id = ? AND author_id = ?", messageID, authorID)
	return err
}
