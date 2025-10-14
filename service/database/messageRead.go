package database

import (
	"time"
)

type MessageRead struct {
	MessageID int64     `json:"message_id"`
	MemberID  int64     `json:"member_id"`
	ReadAt    time.Time `json:"read_at"`
}

// CreateMessageRead inserts or updates a read record for a message
func (db *appdbimpl) CreateMessageRead(messageID, memberID int64) error {
	_, err := db.c.Exec(`
		INSERT OR IGNORE INTO message_read (message_id, member_id, read_at)
		VALUES (?, ?, ?)
	`, messageID, memberID, time.Now().UTC())
	return err
}

// GetReadMembersByMessage returns all member IDs who have read a given message
func (db *appdbimpl) GetReadMembersByMessage(messageID int64) ([]int64, error) {
	rows, err := db.c.Query(`SELECT member_id FROM message_read WHERE message_id = ?`, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []int64
	for rows.Next() {
		var memberID int64
		if err := rows.Scan(&memberID); err != nil {
			return nil, err
		}
		members = append(members, memberID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}
