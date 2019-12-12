package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

type BotConfig struct {
	BotToken  string
	LogDBPath string
	LocalDB   string
}

// LoadBotConfiguration loads the bot options from config.toml
func LoadBotConfiguration() BotConfig {
	return BotConfig{
		BotToken:  os.Getenv("BOT_TOKEN"),
		LogDBPath: os.Getenv("DATABASE_PATH"),
		LocalDB:   os.Getenv("IS_LOCAL"),
	}
}

func main() {
	botConfValue := LoadBotConfiguration()
	bot, err := tgbotapi.NewBotAPI(botConfValue.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	chatLogDB := InitChatDB(botConfValue.LocalDB, botConfValue.LogDBPath)
	log.Printf("Bot authorised on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		// Ignore any non-Message Updates
		if update.Message == nil {
			continue
		}

		im := IncomingMessage{
			MessageID:      update.Message.MessageID,
			MessageTime:    update.Message.Time(),
			FirstName:      update.Message.From.FirstName,
			UserName:       update.Message.From.UserName,
			MessageText:    update.Message.Text,
			SubmittedURL:   "",
			SubmittedImage: nil,
		}

		parsedURL := im.IdentifyMessage()
		if im.IsRepost(parsedURL, chatLogDB) == true {
			im.FlagRepost(*bot, update)
			im.AddReposterToDB(chatLogDB)
		}

		// Add log to DB
		if parsedURL != "" {
			im.AddLogToDB(chatLogDB)
		}
	}
}
