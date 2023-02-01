package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	t "habr-searcher/internal/Tracker"
	"log"
	"os"
	"strconv"
)

type Bot struct {
	TgBot      *tgbotapi.BotAPI
	subChannel chan string
}

func New(sc chan string) *Bot {
	token, exist := os.LookupEnv("TokenForHabrSearcher")
	if !exist {
		log.Fatal("Token for Tg api does not exist")
	}

	tgbot, err := tgbotapi.NewBotAPI(token)
	t.Check(err)

	// tgbot.Debug = true
	return &Bot{
		TgBot:      tgbot,
		subChannel: sc,
	}
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5

	updates := b.TgBot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil /* || !update.Message.IsCommand() */ {
			continue // ignore non-command and non-message
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				// нужно добавить юзера, трекер, и подписать его туда, запрсив тег
				msg.Text = "Hello, bot is active!\nEnter your tag:"
				_, err := b.TgBot.Send(msg)
				t.Check(err)
				//case "addTag":
				//	msg.Text = "Enter new tag:"
				//	_, err := b.TgBot.Send(msg)
				//	t.Check(err)
				//	b.requestNewTag(update.Message.Chat.ID, update.UpdateID+1)
				//	msg.Text = "Your tag successfully added"
				//	_, err = b.TgBot.Send(msg)
				//	t.Check(err)
				//}
			}
		} else {
			id := update.Message.Chat.ID
			tag := update.Message.Text

			tagAndId := fmt.Sprintf("%s %d", tag, id)
			b.subChannel <- tagAndId

			msg.Text = "Your tag successfully added"
			_, err := b.TgBot.Send(msg)
			t.Check(err)
		}
	}

}

func (b *Bot) SendMessage(chatId, text string) {
	idAsNumber, err := strconv.Atoi(chatId)
	t.Check(err)
	msg := tgbotapi.NewMessage(int64(idAsNumber), text)
	_, err = b.TgBot.Send(msg)
	t.Check(err)
}
