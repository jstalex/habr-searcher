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

const startText string = "Hello, I can help you to track new posts from habr.com!\nAs soon as the post is published, I'll send you a link\nSo, enter your tag:"
const helpText string = "Welcome!\n/start - start bot and follow instructions\n/addtag - subscribe on new tag tracking\nFor questions : @jstal3x\nproject in developing..."

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
				msg.Text = startText
				b.send(msg)
			case "addtag":
				msg.Text = "Enter new tag:"
				b.send(msg)
			case "help":
				msg.Text = helpText
				b.send(msg)
			}
		} else {
			id := update.Message.Chat.ID
			tag := update.Message.Text
			if len(tag) > 0 {
				tagAndId := fmt.Sprintf("%s#%d", tag, id)
				b.subChannel <- tagAndId

				msg.Text = "Your tag successfully added"
				b.send(msg)
			} else {
				msg.Text = "Please, try again"
				b.send(msg)
			}

		}
	}

}

func (b *Bot) SendPost(chatId, text string) {
	idAsNumber, err := strconv.Atoi(chatId)
	t.Check(err)
	msg := tgbotapi.NewMessage(int64(idAsNumber), text)
	_, err = b.TgBot.Send(msg)
	t.Check(err)
}

func (b *Bot) send(msg tgbotapi.Chattable) {
	_, err := b.TgBot.Send(msg)
	t.Check(err)
}
