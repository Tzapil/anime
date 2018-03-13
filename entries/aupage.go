package entries

import (
	"strings"

	"golang.org/x/net/html"

	"github.com/tzapil/anime/parse"
	q "github.com/tzapil/anime/query"
)

type Author struct {
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Links []*Link `json:"links"`
	Anime []*Link `json:"anime"`
}

func NewAuthor(n string, e string) *Author {
	return &Author{n, e, make([]*Link, 0), make([]*Link, 0)}
}

func ParseAuthorEmail(node *html.Node) string {
	result := ""
	r := q.NewQuery().
		Child(q.NewTag("tbody")).
		Child(q.NewTag("tr")).
		Child(q.NewTag("td")).
		Child(q.NewTag("table")).
		Child(q.NewTag("tbody")).
		Child(q.NewTag("tr")).
		Child(q.NewTag("td")).
		Find(node)

	if len(r) > 1 {
		fb := r[1]
		if fb != nil && fb.FirstChild != nil {
			result = strings.Replace(fb.FirstChild.Data, "[гав]", "@", -1)
		}
	}

	return result
}

func ParseAuthorName(node *html.Node) string {
	result := ""
	r := q.NewQuery().
		Child(q.NewTag("tbody")).
		Child(q.NewTag("tr")).
		Child(q.NewTag("td")).
		Find(node)

	if len(r) > 2 {
		fb := r[1]
		blq := fb.FirstChild
		if blq != nil && blq.Type == html.TextNode {
			result = blq.Data
		}
	}

	return result
}

func ParseAuthorLinks(node *html.Node) []*Link {
	result := make([]*Link, 0)
	r := q.NewQuery().
		Child(q.NewTag("tbody")).
		Child(q.NewTag("tr")).
		Child(q.NewTag("td")).
		Child(q.NewTag("table")).
		Child(q.NewTag("tbody")).
		Child(q.NewTag("tr")).
		Child(q.NewTag("td")).
		Child(q.NewTag("a")).
		Find(node)

	for i := 0; i < len(r); i++ {
		result = append(result, NewLink(r[i].FirstChild.Data, siteAddress+parse.GetHref(r[i])))
	}

	return result
}

func ParseAuthorWorks(node *html.Node) []*Link {
	result := make([]*Link, 0)
	r := q.NewQuery().
		Child(q.NewTag("tbody")).
		Child(q.NewTag("tr")).
		Child(q.NewTag("td")).
		Child(q.NewTag("blockquote")).
		Child(q.NewTag("a")).
		Find(node)

	for i := 0; i < len(r); i++ {
		fb := r[i]

		if fb != nil {
			result = append(result, NewLink(fb.FirstChild.Data, siteAddress+parse.GetHref(fb)))
		}
	}

	return result
}

func ParseAuthorEntry(node *html.Node) *Author {
	result := NewAuthor(ParseAuthorName(node), ParseAuthorEmail(node))
	result.Anime = ParseAuthorWorks(node)
	result.Links = ParseAuthorLinks(node)
	return result
}

func ParseAuthorPage(root *html.Node) *Author {
	var result *Author = nil

	r := q.NewQuery().
		Child(q.NewTag("html")).
		Child(q.NewTag("body")).
		Child(q.NewTag("table")).
		Child(q.NewTag("tbody")).
		Child(q.NewTag("tr")).
		Child(q.NewTag("td")).
		Child(q.NewTag("table")).
		Find(root)

	if len(r) > 4 {
		// Tables with content
		information := r[4]
		result = ParseAuthorEntry(information)
	}

	return result
}
