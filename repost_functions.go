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
func InitChatDB(IsLocal string, path string) *sql.DB {
	var database *sql.DB
	switch {
	case IsLocal == "true":
		dataPath := filepath.Join(".", "data")
		os.MkdirAll(dataPath, os.ModePerm)
		database, _ = sql.Open("sqlite3", filepath.Join(dataPath, path))
	case IsLocal == "false":
		database, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	}
	chatLogTable :=
		`create table if not exists chatLog(
		message_id integer,
		message_timestamp text,
		username text,
		message_content text,
		submitted_url text
	);`

	repostCountTable :=
		`create table if not exists repostLog(
		first_name text,
		username text primary key,
		repost_count integer	
	);`

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
	insert into chatLog(
		message_id, 
		message_timestamp, 
		username,
		message_content,
		submitted_url
	) values(?, ?, ?, ?, ?)`
	stmt, err := database.Prepare(logEntry)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	stmt.Exec(im.MessageID, im.MessageTime, im.UserName, im.MessageText, im.SubmittedURL)
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

// AddReposterToDB adds the offending reposter to the DB with the number of reposts
func (im *IncomingMessage) AddReposterToDB(database *sql.DB) {
	var repost_count int32
	repostEntry := `
	insert or ignore into repostLog(
		first_name,
		username,
		repost_count) values(?, ?, ?)`
	stmt, err := database.Prepare(`update repostLog set repost_count = repost_count +1 where username = ?`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	stmt.Exec(repostEntry, im.FirstName, im.UserName, repost_count)
}