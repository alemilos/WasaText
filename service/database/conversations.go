package database

import (
	"database/sql"
	"time"
)

// Conversation represents a chat (private or group)
type Conversation struct {
	ID        int64     `json:"conversationId"`
	Type      string    `json:"type"`
	Name      *string   `json:"name,omitempty"`
	PhotoPath *string   `json:"photo_path,omitempty"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// CreatePrivateConversation creates a private conversation between two users
func (db *appdbimpl) CreatePrivateConversation(user1ID, user2ID int64) (*Conversation, error) {
	res, err := db.c.Exec(`
		INSERT INTO conversations (type, created_by) 
		VALUES ('private', ?)`, user1ID)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()

	// Add both users as members
	if err := db.AddMember(id, user1ID, "member"); err != nil {
		return nil, err
	}
	if err := db.AddMember(id, user2ID, "member"); err != nil {
		return nil, err
	}

	return db.GetConversationByID(id)
}

// GetPrivateConversation returns an existing private conversation between two users, if exists
func (db *appdbimpl) GetPrivateConversation(user1ID, user2ID int64) (*Conversation, error) {
	var conv Conversation
	err := db.c.QueryRow(`
		SELECT c.id, c.type, c.name, c.photo_path, c.created_by, c.created_at
		FROM conversations c
		JOIN conversation_members cm1 ON cm1.conversation_id = c.id AND cm1.user_id = ?
		JOIN conversation_members cm2 ON cm2.conversation_id = c.id AND cm2.user_id = ?
		WHERE c.type = 'private' LIMIT 1`,
		user1ID, user2ID).Scan(&conv.ID, &conv.Type, &conv.Name, &conv.PhotoPath, &conv.CreatedBy, &conv.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // no conversation found
	}
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

// GetConversationByID returns a conversation by its ID
func (db *appdbimpl) GetConversationByID(id int64) (*Conversation, error) {
	var conv Conversation
	err := db.c.QueryRow(`
		SELECT id, type, name, photo_path, created_by, created_at
		FROM conversations WHERE id = ?`, id).
		Scan(&conv.ID, &conv.Type, &conv.Name, &conv.PhotoPath, &conv.CreatedBy, &conv.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

// GetConversationsByUserID returns all conversations for a given user
func (db *appdbimpl) GetConversationsByUserID(userID int64) ([]Conversation, error) {
	rows, err := db.c.Query(`
		SELECT c.id, c.type, c.name, c.photo_path, c.created_by, c.created_at
		FROM conversations c
		JOIN conversation_members cm ON cm.conversation_id = c.id
		WHERE cm.user_id = ?
		ORDER BY c.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []Conversation
	for rows.Next() {
		var c Conversation
		if err := rows.Scan(&c.ID, &c.Type, &c.Name, &c.PhotoPath, &c.CreatedBy, &c.CreatedAt); err != nil {
			return nil, err
		}
		convs = append(convs, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return convs, nil
}

// CreateGroupConversation creates a group conversation with name and members
func (db *appdbimpl) CreateGroupConversation(name string, creatorID int64, memberIDs []int64) (*Conversation, error) {
	res, err := db.c.Exec(`
		INSERT INTO conversations (type, name, created_by) 
		VALUES ('group', ?, ?)`, name, creatorID)
	if err != nil {
		return nil, err
	}

	convID, _ := res.LastInsertId()

	// Add creator as admin
	if err := db.AddMember(convID, creatorID, "admin"); err != nil {
		return nil, err
	}

	// Add other members as regular members
	for _, id := range memberIDs {
		if id == creatorID {
			continue // already added as admin
		}
		if err := db.AddMember(convID, id, "member"); err != nil {
			return nil, err
		}
	}

	return db.GetConversationByID(convID)
}
