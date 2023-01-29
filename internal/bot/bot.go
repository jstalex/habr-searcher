package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	t "habr-searcher/internal/Tracker"
	"log"
	"os"
)

type Bot struct {
	TgBot *tgbotapi.BotAPI
}

func New() *Bot {
	token, exist := os.LookupEnv("TokenForHabrSearcher")
	if !exist {
		log.Println("Token for Tg api does not exist")
	}

	tgbot, err := tgbotapi.NewBotAPI(token)
	t.Check(err)
	return &Bot{
		TgBot: tgbot,
	}
}
