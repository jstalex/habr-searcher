package habr_parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"time"
)

type Parser struct {
	SearchReq string
	lastTime  time.Time
}

type Post struct {
	Name   string
	Link   string
	Author string
}

func New(req string) *Parser {
	return &Parser{
		SearchReq: req,
		lastTime:  time.Now(),
	}
}

func (p *Parser) GetNewPost() (Post, bool) {
	newPost := Post{}

	url := fmt.Sprintf("https://habr.com/ru/search/?q=%s&target_type=posts&order=date", p.SearchReq)
	res, err := http.Get(url)
	check(err)
	body := res.Body
	defer body.Close()

	//io.Copy(os.Stdout, searchPage)

	searchPage, err := goquery.NewDocumentFromReader(body)
	check(err)

	newestPost := searchPage.Find(".tm-articles-list").Find("article").First()

	publishTime, ok := newestPost.Find(".tm-article-snippet__datetime-published").Find("time").Attr("datetime")
	attrIsOk(ok)

	pTime, err := time.Parse(time.RFC3339, publishTime)
	check(err)

	if pTime.After(p.lastTime) {
		return newPost, false
	}

	postId, ok := newestPost.Attr("id")
	attrIsOk(ok)

	author, ok := newestPost.Find(".tm-user-info tm-article-snippet__author").Find("a").Attr("title")
	attrIsOk(ok)

	postName, err := newestPost.Find("h2").Find("a").Find("span").Html()
	check(err)

	link, ok := newestPost.Find("h2").Find("a").Attr("href")
	attrIsOk(ok)
	link = fmt.Sprintf("https://habr.com%s", link)

	fmt.Println(postId, publishTime, postName, author, link, pTime)
	return newPost, true
}

func check(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func attrIsOk(ok bool) {
	if !ok {
		log.Println("Such Attribute not exist")
	}
}
