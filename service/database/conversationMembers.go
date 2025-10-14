package database

import (
	"fmt"
	"time"
)

// ConversationMember represents a user in a conversation
type ConversationMember struct {
	ConversationID int64     `json:"conversationId"`
	UserID         int64     `json:"userId"`
	Role           *string   `json:"role,omitempty"`
	JoinedAt       time.Time `json:"joined_at"`
}

// AddMember adds a user to a conversation
func (db *appdbimpl) AddMember(conversationID, userID int64, role string) error {
	_, err := db.c.Exec(`
		INSERT INTO conversation_members (conversation_id, user_id, role) 
		VALUES (?, ?, ?)`, conversationID, userID, role)
	return err
}

// RemoveMember removes a user from a conversation
func (db *appdbimpl) RemoveMember(conversationID, userID int64) error {
	_, err := db.c.Exec(`
		DELETE FROM conversation_members 
		WHERE conversation_id = ? AND user_id = ?`, conversationID, userID)
	return err
}

func (db *appdbimpl) IsMember(conversationID, userID int64) error {
	var exists bool
	err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM conversation_members 
			WHERE conversation_id = ? AND user_id = ?
		)
	`, conversationID, userID).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to check membership: %w", err)
	}

	if !exists {
		return fmt.Errorf("user %d is not a member of conversation %d", userID, conversationID)
	}

	return nil
}

func (db *appdbimpl) GetMembersByConversation(conversationID int64) ([]int64, error) {
	rows, err := db.c.Query(`SELECT user_id FROM conversation_members WHERE conversation_id = ?`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		members = append(members, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}
