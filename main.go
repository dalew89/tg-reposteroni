package main

import (
	"github.com/BurntSushi/toml"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type BotConfigFromFile struct {
	BotToken    string `toml:"bot_token"`
	LogDBPath   string `toml:"repost_log_file"`
	ToggleDebug bool   `toml:"debug_enabled"`
}

// LoadBotConfiguration loads the bot options from config.toml
func LoadBotConfiguration() BotConfigFromFile {
	var config BotConfigFromFile
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}
	return BotConfigFromFile{
		BotToken:    config.BotToken,
		LogDBPath:   config.LogDBPath,
		ToggleDebug: config.ToggleDebug,
	}
}

func main() {
	botConf := LoadBotConfiguration()
	bot, err := tgbotapi.NewBotAPI(botConf.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	chatLogDB := InitChatDB(botConf.LogDBPath)
	log.Printf("Auth'd on account %s", bot.Self.UserName)
	bot.Debug = botConf.ToggleDebug
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
		im := IncomingMessage{
			MessageID:      update.Message.MessageID,
			MessageTime:    update.Message.Time(),
			UserName:       update.Message.From.UserName,
			MessageText:    update.Message.Text,
			SubmittedURL:   "",
			SubmittedImage: nil,
		}
		url := im.IdentifyMessage()
		if im.IsRepost(url, chatLogDB) == true {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				`╰( ͡° ͜ʖ ͡° )つ──☆*:・ﾟ Copypastus Totalus!! I can't believe people actually take time out
of their day to copy and paste links instead of contributing to chat.`)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
		if url != "" {
			im.AddLogToDB(chatLogDB)
		}
	}
}
