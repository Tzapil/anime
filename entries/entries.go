package entries

import (
	"github.com/tzapil/anime/parse"
	q "github.com/tzapil/anime/query"
	"golang.org/x/net/html"
)

const siteAddress = "http://www.fansubs.ru/"

type Link struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

func NewLink(n string, l string) *Link {
	return &Link{n, l}
}

func get_date(head *html.Node) string {
	result := ""
	r := q.NewQuery().
		Child(q.NewTag("tbody")).
		Child(q.NewTag("tr")).
		Child(q.NewTag("td")).
		Find(head)

	if len(r) >= 5 {
		fb := r[4].FirstChild
		if fb != nil && fb.Type == html.TextNode {
			result = fb.Data
		}
	}

	return result
}

func get_info(head *html.Node) string {
	result := ""
	r := q.NewQuery().
		Child(q.NewTag("tbody")).
		Child(q.NewTag("tr")).
		Child(q.NewTag("td")).
		Child(q.NewTag("b")).
		Find(head)

	if len(r) > 0 {
		fb := r[0].FirstChild
		if fb != nil && fb.Type == html.TextNode {
			result = fb.Data
		}
	}

	return result
}

func get_format(head *html.Node) string {
	result := ""
	r := q.NewQuery().
		Child(q.NewTag("tbody")).
		Child(q.NewTag("tr")).
		Child(q.NewTag("td")).
		Child(q.NewTag("a")).
		Child(q.NewTag("font")).
		Find(head)

	if len(r) > 0 {
		fb := r[0].FirstChild
		if fb != nil && fb.Type == html.TextNode {
			result = fb.Data
		}
	}

	return result
}

func get_download(head *html.Node) string {
	result := ""
	base := ""
	srt := ""

	r := q.NewQuery().Child(q.NewTag("form")).Find(head)

	if len(r) > 0 {
		base = parse.GetAttr(r[0], "action")
	}

	r = q.NewQuery().Child(q.NewTag("input")).Find(head)

	if len(r) > 0 {
		srt = parse.GetAttr(r[0], "value")
	}

	if base != "" && srt != "" {
		result = siteAddress + base + "?srt=" + srt
	}

	return result
}
