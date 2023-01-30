package app

import (
	t "habr-searcher/internal/Tracker"
	"habr-searcher/internal/bot"
	"log"
	"strings"
)

type App struct {
	Trackers    map[string]*t.Tracker
	TgBot       *bot.Bot
	UsersForTag map[string][]User
	users       []User
	subChannel  chan string
}

type User struct {
	chatId string
}

func New() *App {
	Trackers := make(map[string]*t.Tracker)
	UsersForTag := make(map[string][]User)
	users := make([]User, 0)
	sc := make(chan string)
	tgBot := bot.New(sc)

	return &App{
		Trackers:    Trackers,
		TgBot:       tgBot,
		UsersForTag: UsersForTag,
		users:       users,
		subChannel:  sc,
	}
}

func (a *App) AddNewTracker(tag string) {
	tracker := t.New(tag)
	a.Trackers[tag] = tracker
}

func (a *App) Run() {
	//a.AddNewTracker("go")
	//post, _ := a.Trackers["go"].GetNewPost()
	//fmt.Printf("%v\n", post)
	go a.TgBot.Run()
	for {
		a.CheckNewSubscribe()
	}
}

func (a *App) AddNewUser(u User) {
	a.users = append(a.users, u)
}

func (a *App) SubscribeNewTagToUser(u User, tag string) {
	if _, ok := a.Trackers[tag]; !ok {
		a.AddNewTracker(tag)
		a.UsersForTag[tag] = make([]User, 0)
		a.UsersForTag[tag] = append(a.UsersForTag[tag], u)
		log.Println("SUCCESS")
		log.Println(a.Trackers)
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

func (a *App) CheckNewSubscribe() {
	str, ok := <-a.subChannel
	values := strings.Split(str, " ")
	tag, id := values[0], values[1]
	if ok {
		a.SubscribeNewTagToUser(User{id}, tag)
	}
}
