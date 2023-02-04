package app

import (
	"context"
	t "habr-searcher/internal/Tracker"
	"habr-searcher/internal/bot"
	"habr-searcher/internal/model"
	"habr-searcher/internal/store"
	"log"
	"strings"
	"time"
)

const updateDelay int = 3 // in seconds

type App struct {
	Trackers     map[string]*t.Tracker
	TgBot        *bot.Bot
	Store        *store.Store
	UsersForTag  map[string][]Chat
	subChannel   chan string
	storeChannel chan model.User
}

type Chat struct {
	chatId string
}

func New() *App {
	Trackers := make(map[string]*t.Tracker)
	UsersForTag := make(map[string][]Chat)
	sc := make(chan string)
	storeChan := make(chan model.User)
	tgBot := bot.New(sc, storeChan)

	ctx := context.Background()
	s := store.New(ctx)

	return &App{
		Trackers:     Trackers,
		TgBot:        tgBot,
		Store:        s,
		UsersForTag:  UsersForTag,
		subChannel:   sc,
		storeChannel: storeChan,
	}
}

func (a *App) AddNewTracker(tag string) {
	tracker := t.New(tag)
	a.Trackers[tag] = tracker
}

func (a *App) Run() {
	go a.TgBot.Run()
	go a.CheckNewSubscribe()
	go a.CheckNewUpdateToStore()
	a.CheckNewPosts()
}

func (a *App) SubscribeNewTagToUser(u Chat, tag string) {
	if _, ok := a.Trackers[tag]; !ok {
		a.AddNewTracker(tag)
		a.UsersForTag[tag] = make([]Chat, 0)
		a.UsersForTag[tag] = append(a.UsersForTag[tag], u)
		//log.Println("SUCCESS")
		//log.Println(a.Trackers)
	} else {
		a.UsersForTag[tag] = append(a.UsersForTag[tag], u)
		//log.Println("SUCCESS")
		//log.Println(a.Trackers)
	}
}

func (a *App) CheckNewPosts() {
	for {
		for tag, tracker := range a.Trackers {
			post, exist := tracker.GetNewPost()
			if exist {
				for _, user := range a.UsersForTag[tag] {
					a.TgBot.SendPost(user.chatId, post.InString())
				}
			}
		}
		time.Sleep(time.Duration(updateDelay) * time.Second)
	}
}

func (a *App) CheckNewSubscribe() {
	// auto locking by subChannel
	for {
		str, ok := <-a.subChannel
		log.Printf("im got info from chan %s %v\n", str, ok)
		values := strings.Split(str, "#")
		tag, id := values[0], values[1]
		if ok {
			a.SubscribeNewTagToUser(Chat{id}, tag)
		} else {
			log.Println("read from subChannel error")
		}
	}
}

func (a *App) CheckNewUpdateToStore() {
	for {
		user, ok := <-a.storeChannel
		if ok {
			err := a.Store.Set(user.UserName, user)
			t.Check(err)
		} else {
			log.Println("read from storeChannel error")
		}
	}
}
