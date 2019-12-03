package main

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mvdan/xurls"
	"time"
)

type IncomingMessage struct {
	MessageID      int
	MessageTime    time.Time
	UserName       string
	MessageText    string
	SubmittedURL   string
	SubmittedImage []interface{}
}

// IdentifyMessage identifies and returns what a message contains
func IdentifyMessage(update tgbotapi.Update) string {
	im := IncomingMessage{
		MessageID:      update.Message.MessageID,
		MessageTime:    update.Message.Time(),
		MessageText:    update.Message.Text,
		SubmittedURL:   "",
		SubmittedImage: nil,
	}
	potentialURL := FindURLInText(im.MessageText)
	switch {
	case potentialURL != "": // if there is a URl inside of the message
		im.SubmittedURL = potentialURL
		return im.SubmittedURL

	case potentialURL == "": // if there isn't

	}
	return im.SubmittedURL
}

// FindURLInText will parse a URL from any text and return a string
func FindURLInText(message string) string {
	rxRelaxed := xurls.Relaxed()
	foundURL := rxRelaxed.FindString(message)
	return foundURL
}

// InitChatDB initialises the database to store the chat logs from the group chat.
func InitChatDB(path string) *sql.DB {
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

// AddLogToDB writes a single chat logs to the the DB.
func (im *IncomingMessage) AddLogToDB(database *sql.DB) {
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

func CheckForRepost(potentialRepostedURL string, database *sql.DB) {
	rows, _ := database.Query(`SELECT submitted_url, message_id, FROM chatLog`)
	for rows.Next() {
		rows.Scan(potentialRepostedURL)
		fmt.Println(potentialRepostedURL)
	}
}
