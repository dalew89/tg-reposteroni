package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

//AddLogToDB adds lines of the group chat to the db file specified in config.toml
func (im *IncomingMessage) AddLogToDB(logFile string) {
	database, _ := sql.Open("sqlite3", logFile)
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS chatLog (id INTEGER PRIMARY KEY, username TEXT, messagetext TEXT, message_time_sent TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO chatLog (username, messagetext, message_time_sent) VALUES (?, ?, ?)")
	statement.Exec(im.UserName, im.MessageText, im.MessageTime)
	//rows, _ := database.Query("SELECT id, username, messagetext, FROM chathistory")
	//var id int
	//var username string
	//var messagetext string
	//for rows.Next() {
	//	rows.Scan(&id, &username, &messagetext)
	//	fmt.Println(strconv.Itoa(id) + ": " + username + " " + messagetext + " ")
	//
	//}
}
