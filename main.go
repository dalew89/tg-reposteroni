package main

import (
	"github.com/BurntSushi/toml"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mvdan/xurls"
	"log"
)

type BotConfigFromFile struct {
	BotToken    string `toml:"bot_token"`
	LogFilePath string `toml:"repost_log_file"`
	ToggleDebug bool   `toml:"debug_enabled"`
}

func LoadBotConfiguration() BotConfigFromFile {
	var config BotConfigFromFile
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}
	return BotConfigFromFile{BotToken: config.BotToken, LogFilePath: config.LogFilePath}
}

func main() {
	botConf := LoadBotConfiguration()
	bot, err := tgbotapi.NewBotAPI(botConf.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = botConf.ToggleDebug
	db := InitChatLogDB(botConf.LogFilePath)
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
			MessageID:    update.Message.MessageID,
			MessageTime:  update.Message.Time(),
			UserName:     update.Message.From.UserName,
			MessageText:  update.Message.Text,
			SubmittedURL: "",
		}

		rxRelaxed := xurls.Relaxed()
		foundURL := rxRelaxed.FindString(update.Message.Text)
		if foundURL != "" { // If the group message contains a URL, save it to the database
			im.SubmittedURL = foundURL
			im.StoreChatLog(db)
		}
	}
}
