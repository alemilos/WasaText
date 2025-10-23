package database

import (
	"database/sql"
	"time"
)

// Message represents a chat message
type Message struct {
	ID             int64     `json:"id"`
	Type           string    `json:"type"` // "text" or "image"
	IsForwarded    bool      `json:"is_forwarded"`
	ConversationID int64     `json:"conversation_id"`
	AuthorID       int64     `json:"author_id"`
	Content        string    `json:"content"`
	SecondaryContent string `json:"secondary_content"`
	CreatedAt      time.Time `json:"created_at"`
}

// CreateMessage inserts a new message into a conversation
func (db *appdbimpl) CreateMessage(conversationID, authorID int64, msgType string, content string, secondaryContent string, isForwarded bool) (*Message, error) {
	res, err := db.c.Exec(`
		INSERT INTO messages (type, is_forwarded, conversation_id, author_id, content, secondary_content)
		VALUES (?, ?, ?, ?, ?)`,
		msgType, boolToInt(isForwarded), conversationID, authorID, content, secondaryContent)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()
	return db.GetMessageByID(id)
}

// GetMessageByID fetches a message by ID
func (db *appdbimpl) GetMessageByID(id int64) (*Message, error) {
	var m Message
	var isForwardedInt int
	err := db.c.QueryRow(`
		SELECT id, type, is_forwarded, conversation_id, author_id, content, created_at
		FROM messages WHERE id = ?`, id).
		Scan(&m.ID, &m.Type, &isForwardedInt, &m.ConversationID, &m.AuthorID, &m.Content, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	m.IsForwarded = intToBool(isForwardedInt)
	return &m, nil
}

// GetMessagesByConversation fetches all messages in a conversation
func (db *appdbimpl) GetMessagesByConversation(conversationID int64) ([]Message, error) {
	rows, err := db.c.Query(`
		SELECT id, type, is_forwarded, conversation_id, author_id, content, created_at
		FROM messages
		WHERE conversation_id = ?
		ORDER BY created_at ASC`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		var isForwardedInt int
		if err := rows.Scan(&m.ID, &m.Type, &isForwardedInt, &m.ConversationID, &m.AuthorID, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		m.IsForwarded = intToBool(isForwardedInt)
		msgs = append(msgs, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return msgs, nil
}

// UpdateMessageContent updates the content of an existing message (used for images after saving file)
// This shouldn't be used to "EDIT" the message since this functionality is not provided by the API.
func (db *appdbimpl) UpdateMessageContent(messageID int64, content string) error {
	_, err := db.c.Exec(`UPDATE messages SET content = ? WHERE id = ?`, content, messageID)
	return err
}

// GetLastMessageByConversation fetches the latest message for a conversation
func (db *appdbimpl) GetLastMessageByConversation(conversationID int64) (*Message, error) {
	var m Message
	var isForwardedInt int
	err := db.c.QueryRow(`
		SELECT id, type, is_forwarded, conversation_id, author_id, content, created_at
		FROM messages
		WHERE conversation_id = ?
		ORDER BY created_at DESC
		LIMIT 1`, conversationID).
		Scan(&m.ID, &m.Type, &isForwardedInt, &m.ConversationID, &m.AuthorID, &m.Content, &m.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	m.IsForwarded = intToBool(isForwardedInt)
	return &m, nil
}

func (db *appdbimpl) DeleteMessage(messageID int64) error {
	res, err := db.c.Exec(`DELETE FROM messages WHERE id = ?`, messageID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return sql.ErrNoRows
	}
	return err
}

// Helper to convert bool <-> int for SQLite
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func intToBool(i int) bool {
	return i != 0
}
