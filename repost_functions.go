package main

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"mvdan.cc/xurls"
	"os"
	"path/filepath"
	"time"
)

type IncomingMessage struct {
	MessageID      int
	MessageTime    time.Time
	FirstName      string
	LastName       string
	UserName       string
	MessageText    string
	SubmittedURL   string
	SubmittedImage []interface{}
}

// IdentifyMessage identifies and returns what a message contains
func (im *IncomingMessage) IdentifyMessage() string {
	potentialURL := FindURLInText(im.MessageText)
	switch {
	case potentialURL != "": // if there is a URl inside of the message
		im.SubmittedURL = potentialURL
		return im.SubmittedURL
	default: // if there isn't
		return ""
	}
}

// FindURLInText will parse a URL from any text and return a string
func FindURLInText(message string) string {
	foundURL := xurls.Relaxed().FindString(message)
	return foundURL
}

// InitChatDB initialises the database to store the chat logs from the group chat.
func InitChatDB(databasePath string) *sql.DB {
	dataPath := filepath.Join(".", "data")
	os.MkdirAll(dataPath, os.ModePerm)
	database, _ := sql.Open("sqlite3", filepath.Join(dataPath, databasePath))

	chatLogTable := `create table if not exists chatLog (
		message_id integer,
		message_timestamp text,
		first_name text,
		last_name text,
		username text,
		submitted_url text);`

	repostCountTable := `create table if not exists repostLog (
		first_name text,
		last_name text,
		username text primary key,
		repost_count integer);`

	_, err := database.Exec(chatLogTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err2 := database.Exec(repostCountTable)
	if err2 != nil {
		log.Fatal(err2)
	}
	return database
}

// AddLogToDB writes a single chat logs to the the DB.
func (im *IncomingMessage) AddLogToDB(database *sql.DB) {
	logEntry := `
	insert into chatLog (
		message_id, 
		message_timestamp,
		first_name,
		last_name,  
		username,
		submitted_url) 
		values(?, ?, ?, ?, ?, ?)`
	statement, err := database.Prepare(logEntry)
	if err != nil {
		panic(err)
	}
	defer statement.Close()
	statement.Exec(im.MessageID,
		im.MessageTime,
		im.FirstName,
		im.LastName,
		im.UserName,
		im.SubmittedURL)
}

// AddReposterToDB adds the offending reposter to the DB with the number of reposts
func (im *IncomingMessage) AddReposterToDB(database *sql.DB) {
	repostLogEntry := `insert into repostLog (
		first_name,
		last_name,
		username,
		repost_count) values(?, ?, ?, 0)`

	addRow, err := database.Prepare(repostLogEntry)
	if err != nil {
		log.Fatal(err)
	}
	defer addRow.Close()
	addRow.Exec(im.FirstName, im.LastName, im.UserName)

	updateCount := `update repostLog 
		set repost_count = repost_count+1 
		where username = ? 
		or first_name = ? and last_name = ?`

	incrementRow, err2 := database.Prepare(updateCount)
	if err2 != nil {
		log.Fatal(err)
	}
	defer incrementRow.Close()
	incrementRow.Exec(im.UserName, im.FirstName, im.LastName)
}

//RetrieveRepostStats queries the database for a list of all reposters in the chat
func RetrieveRepostStats(database *sql.DB) {
	// TODO
}

// IsRepost scans the db for potential URL reposts. If it is a repost, return true
func (im *IncomingMessage) IsRepost(potentialRepostedURL string, database *sql.DB) bool {
	rows, err := database.Query(
		`select submitted_url from chatLog where submitted_url = ?`, potentialRepostedURL)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var messageID int
	numberOfLinks := 0
	for rows.Next() {
		rows.Scan(&potentialRepostedURL, &messageID)
		numberOfLinks += 1
	}
	switch {
	case numberOfLinks > 0:
		return true
	default:
		return false
	}
}

// FlagRepost replies to the reposted link
func (im *IncomingMessage) FlagRepost(bot tgbotapi.BotAPI, update tgbotapi.Update) {
	repostWarning := tgbotapi.NewMessage(update.Message.Chat.ID,
		"╰( ͡° ͜ʖ ͡° )つ──☆*:・ﾟ \n"+
			"Repostus Copypastus Totalus!!\n"+
			"I can't believe people actually take time out of their day to copy and paste links "+
			"instead of contributing to chat.")
	repostWarning.ReplyToMessageID = update.Message.MessageID
	bot.Send(repostWarning)
}
