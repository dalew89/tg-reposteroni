package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

type BotConfig struct {
	BotToken    string
	LogDBPath   string
	LocalDB     string
}

// LoadBotConfiguration loads the bot options from config.toml
func LoadBotConfiguration() BotConfig {
	return BotConfig{
		BotToken:    os.Getenv("BOT_TOKEN"),
		LogDBPath:   os.Getenv("DATABASE_PATH"),
		LocalDB:     os.Getenv("IS_LOCAL"),
	}
}

func main() {
	botConf := LoadBotConfiguration()
	bot, err := tgbotapi.NewBotAPI(botConf.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	chatLogDB := InitChatDB(botConf.LocalDB, botConf.LogDBPath)
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
		im := IncomingMessage{
			MessageID:      update.Message.MessageID,
			MessageTime:    update.Message.Time(),
			UserName:       update.Message.From.UserName,
			MessageText:    update.Message.Text,
			SubmittedURL:   "",
			SubmittedImage: nil,
		}
		parsedURL := im.IdentifyMessage()
		if im.IsRepost(parsedURL, chatLogDB) == true {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"╰( ͡° ͜ʖ ͡° )つ──☆*:・ﾟ \n"+
					"Repostus Copypastus Totalus!!\n"+
					"I can't believe people actually take time out of their day to copy and paste links instead"+
					" of contributing to chat.")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
		if parsedURL != "" {
			im.AddLogToDB(chatLogDB)
		}
	}
}
