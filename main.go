package main

import (
	"github.com/BurntSushi/toml"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"time"
)

type BotConfigFromFile struct {
	BotToken string `toml:"bot_token"`
	LogFile  string `toml:"repost_log_file"`
}

type IncomingMessage struct {
	MessageID   int
	MessageTime time.Time
	UserName    string
	MessageText string
}

func loadBotConfig() BotConfigFromFile {
	var config BotConfigFromFile
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}
	return BotConfigFromFile{BotToken: config.BotToken, LogFile: config.LogFile}
}

func main() {
	botConf := loadBotConfig()
	bot, err := tgbotapi.NewBotAPI(botConf.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	log.Printf("Auth'd on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		im := &IncomingMessage{
			MessageID:   update.Message.MessageID,
			MessageTime: update.Message.Time(),
			UserName:    update.Message.From.UserName,
			MessageText: update.Message.Text,
		}

		im.AddLogToDB(botConf.LogFile)

		//log.Printf("Message text: %s", im.MessageText)
		//log.Printf("Message received from: %s", im.UserName)

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID
		//bot.Send(msg)
	}
}
