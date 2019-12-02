package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type IncomingMessage struct {
	MessageID    int
	MessageTime  time.Time
	UserName     string
	MessageText  string
	SubmittedURL string
}

// InitChatLogDB initialises the DB to store the chat logs from the group chat.
func InitChatLogDB(path string) *sql.DB {
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	if database == nil {
		panic(err)
	}
	chatLogTable := `
	CREATE TABLE IF NOT EXISTS chatLog(
		message_id INTEGER,
		message_timestamp TEXT,
		username TEXT,
		message_content TEXT,
		submitted_url TEXT
	);
	`
	database.Exec(chatLogTable)
	return database
}

// StoreChatLog writes chat logs to the DB
func (im *IncomingMessage) StoreChatLog(database *sql.DB) {
	logEntry := `
	INSERT INTO chatLog(
		message_id, 
		message_timestamp, 
		username,
		message_content,
		submitted_url
	) values(?, ?, ?, ?, ?)
	`
	stmt, err := database.Prepare(logEntry)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	stmt.Exec(im.MessageID, im.MessageTime, im.UserName, im.MessageText, im.SubmittedURL)
}
