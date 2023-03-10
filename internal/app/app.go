package app

import (
	t "habr-searcher/internal/Tracker"
	"habr-searcher/internal/bot"
	"log"
	"strings"
	"time"
)

const updateDelay int = 3 // in seconds

type App struct {
	Trackers    map[string]*t.Tracker
	TgBot       *bot.Bot
	UsersForTag map[string][]User
	users       []User
	subChannel  chan string
	//wg          *sync.WaitGroup
}

type User struct {
	chatId string
}

func New() *App {
	Trackers := make(map[string]*t.Tracker)
	UsersForTag := make(map[string][]User)
	sc := make(chan string)
	tgBot := bot.New(sc)
	//wg := &sync.WaitGroup{}

	return &App{
		Trackers:    Trackers,
		TgBot:       tgBot,
		UsersForTag: UsersForTag,
		subChannel:  sc,
		//wg:          wg,
	}
}

func (a *App) AddNewTracker(tag string) {
	tracker := t.New(tag)
	a.Trackers[tag] = tracker
}

func (a *App) Run() {
	go a.TgBot.Run()
	go a.CheckNewSubscribe()
	a.CheckNewPosts()
}

func (a *App) SubscribeNewTagToUser(u User, tag string) {
	if _, ok := a.Trackers[tag]; !ok {
		a.AddNewTracker(tag)
		a.UsersForTag[tag] = make([]User, 0)
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
	//defer a.wg.Done()
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
	//defer a.wg.Done()
	// auto locking by subChannel
	for {
		str, ok := <-a.subChannel
		log.Printf("im got info from chan %s %v\n", str, ok)
		values := strings.Split(str, "#")
		tag, id := values[0], values[1]
		if ok {
			a.SubscribeNewTagToUser(User{id}, tag)
		}
	}
}
