package Tracker

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"time"
)

type Tracker struct {
	SearchReq string
	lastTime  time.Time
}

func New(req string) *Tracker {
	t, err := time.Parse(time.RFC3339, "2022-02-22T10:00:00.000Z")
	Check(err)
	return &Tracker{
		SearchReq: req,
		lastTime:  t,
	}
}

func (p *Tracker) GetNewPost() (Post, bool) {
	newPost := Post{}

	url := fmt.Sprintf("https://habr.com/ru/search/?q=%s&target_type=posts&order=date", p.SearchReq)
	res, err := http.Get(url)
	log.Printf("habr req status = %s\n", res.Status)
	Check(err)

	body := res.Body
	defer body.Close()

	//io.Copy(os.Stdout, searchPage)

	searchPage, err := goquery.NewDocumentFromReader(body)
	Check(err)

	newestPost := searchPage.Find(".tm-articles-list").Find("article").First()

	if len(newestPost.Text()) == 0 {
		log.Println("This tag's posts not exist")
		return newPost, false
	}

	publishTime, ok := newestPost.Find(".tm-article-snippet__datetime-published").Find("time").Attr("datetime")
	attrIsOk(ok)

	pTime, err := time.Parse(time.RFC3339, publishTime)
	Check(err)

	if !pTime.After(p.lastTime) {
		return newPost, false
	}

	//postId, ok := newestPost.Attr("id")
	//attrIsOk(ok)

	author, ok := newestPost.Find(".tm-user-info__userpic").Attr("title")
	attrIsOk(ok)

	postName, err := newestPost.Find("h2").Find("a").Find("span").Html()
	Check(err)

	link, ok := newestPost.Find("h2").Find("a").Attr("href")
	attrIsOk(ok)
	link = fmt.Sprintf("https://habr.com%s", link)

	newPost.Name = postName
	newPost.Link = link
	newPost.Author = author

	p.lastTime = pTime

	return newPost, true
}

func Check(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func attrIsOk(ok bool) {
	if !ok {
		log.Println("Such Attribute not exist")
	}
}
