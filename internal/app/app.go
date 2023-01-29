package app

import (
	"fmt"
	t "habr-searcher/internal/Tracker"
	"habr-searcher/internal/bot"
)

type App struct {
	Trackers map[string]*t.Tracker
	TgBot    *bot.Bot
}

func New() *App {
	Trackers := make(map[string]*t.Tracker)
	tgBot := bot.New()
	return &App{
		Trackers: Trackers,
		TgBot:    tgBot,
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
