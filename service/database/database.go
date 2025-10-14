/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// Users
	CreateUser(username string, photoPath *string) (*User, error)
	SetUsername(userID int64, newName string) error
	SetProfilePhoto(userID int64, photoPath string) error
	GetUserById(id int64) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUsers() ([]User, error)

	// Conversations
	CreatePrivateConversation(user1ID, user2ID int64) (*Conversation, error)
	GetPrivateConversation(user1ID, user2ID int64) (*Conversation, error)
	GetConversationByID(id int64) (*Conversation, error)
	GetConversationsByUserID(userID int64) ([]Conversation, error)
	CreateGroupConversation(name string, creatorID int64, memberIDs []int64) (*Conversation, error)

	// Conversation Members
	IsMember(conversationID, userID int64) error
	GetMembersByConversation(conversationID int64) ([]int64, error)

	// Messages
	CreateMessage(conversationID, authorID int64, msgType string, content string, isForwarded bool) (*Message, error)
	GetMessageByID(id int64) (*Message, error)
	GetMessagesByConversation(conversationID int64) ([]Message, error)
	GetLastMessageByConversation(conversationID int64) (*Message, error)
	UpdateMessageContent(messageID int64, content string) error
	DeleteMessage(id int64) error

	// Message Read
	CreateMessageRead(messageID, memberID int64) error
	GetReadMembersByMessage(messageID int64) ([]int64, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY, name TEXT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	// chiama la funzione che crea le tabelle se mancano
	if err := InitializeSchema(db); err != nil {
		return nil, fmt.Errorf("initialize schema: %w", err)
	}

	return &appdbimpl{c: db}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
