package app

import (
	"fmt"
	t "habr-searcher/internal/Tracker"
)

type App struct {
	Trackers map[string]*t.Tracker
}

func New() *App {
	Trackers := make(map[string]*t.Tracker)
	return &App{
		Trackers: Trackers,
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
