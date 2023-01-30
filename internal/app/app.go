package app

import (
	"fmt"
	t "habr-searcher/internal/Tracker"
	"habr-searcher/internal/bot"
)

type App struct {
	Trackers    map[string]*t.Tracker
	TgBot       *bot.Bot
	UsersForTag map[string][]User
	users       []User
}

type User struct {
	chatId string
}

func New() *App {
	Trackers := make(map[string]*t.Tracker)
	UsersForTag := make(map[string][]User)
	users := make([]User, 0)
	tgBot := bot.New()

	return &App{
		Trackers:    Trackers,
		TgBot:       tgBot,
		UsersForTag: UsersForTag,
		users:       users,
	}
}

func (a *App) AddNewTracker(tag string) {
	tracker := t.New(tag)
	a.Trackers[tag] = tracker
}

func (a *App) Run() {
	a.AddNewTracker("go")
	post, _ := a.Trackers["go"].GetNewPost()
	fmt.Printf("%v\n", post)
}

func (a *App) AddNewUser(u User) {
	a.users = append(a.users, u)
}

func (a *App) SubscribeNewTagToUser(u User, tag string) {
	if _, ok := a.Trackers[tag]; !ok {
		a.AddNewTracker(tag)
		a.UsersForTag[tag] = make([]User, 0)
		a.UsersForTag[tag] = append(a.UsersForTag[tag], u)
	} else {
		a.UsersForTag[tag] = append(a.UsersForTag[tag], u)
	}
}

//func (a *App) CheckNewPosts() {
//	for tag, tracker := range a.Trackers{
//		post, exist := tracker.GetNewPost()
//		if exist {
//			for _, user := range a.UsersForTag[tag] {
//				// send message from bot
//			}
//		}
//	}
//}
