package Tracker

import "fmt"

type Post struct {
	Name   string
	Link   string
	Author string
}

func (p *Post) InString() string {
	return fmt.Sprintf("Author: %s\n%s\n%s\n", p.Author, p.Name, p.Link)
}
