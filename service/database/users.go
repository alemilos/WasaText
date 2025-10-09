package database

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64      `json:"id"`
	Username  string     `json:"username"`
	PhotoPath *string    `json:"photo_path,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// Create a user  
func (db *appdbimpl) CreateUser(username string, photoPath *string) (*User, error) {
	res, err := db.c.Exec(`INSERT INTO users (username, photo_path) VALUES (?, ?)`, username, photoPath)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()
	return db.GetUser(id)
}

// Set a new username for a given user, by user id
func (db *appdbimpl) SetUsername(userID int64, newName string) error {
	_, err := db.c.Exec(`UPDATE users SET username=? WHERE id=?`, newName, userID)
	return err
}

// Set a user's profile photo path for a given user, by user id
func (db *appdbimpl) SetProfilePhoto(userID int64, photoPath string) error {
	_, err := db.c.Exec(`UPDATE users SET photo_path=? WHERE id=?`, photoPath, userID)
	return err
}

// Get a user by id
func (db *appdbimpl) GetUser(id int64) (*User, error) {
	var u User
	err := db.c.QueryRow(`SELECT id, username, photo_path, created_at FROM users WHERE id=?`, id).
		Scan(&u.ID, &u.Username, &u.PhotoPath, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Get a user by username
func (db *appdbimpl) GetUserByUsername(username string) (*User, error) {
	var u User
	err := db.c.QueryRow(`
		SELECT id, username, photo_path, created_at 
		FROM users 
		WHERE username = ?`, username,
	).Scan(&u.ID, &u.Username, &u.PhotoPath, &u.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // means not found
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Get all the users in the database 
func (db *appdbimpl) GetUsers() ([]User, error) {
	rows, err := db.c.Query(`SELECT id, username, photo_path, created_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.PhotoPath, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
