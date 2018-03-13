package query

import (
	"golang.org/x/net/html"
)

type Query struct {
	query []Selector
}

func NewQuery() *Query {
	return &Query{make([]Selector, 0)}
}

func (q *Query) Tag(name string) *Query {
	q.query = append(q.query, NewTag(name))

	return q
}

func (q *Query) Child(s Selector) *Query {
	q.query = append(q.query, NewChild(s))
	return q
}

func (q *Query) Next() *Query {
	q.query = append(q.query, NewNext())
	return q
}

func (q *Query) Contains(val string) *Query {
	q.query = append(q.query, NewContains(val))
	return q
}

func recursion(n *html.Node, s Selector, r *[]*html.Node) {
	if n != nil {
		if s.Check(n) {
			*r = append(*r, n)
		}
		recursion(n.FirstChild, s, r)
		recursion(n.NextSibling, s, r)
	}
}

func GetBySelector(n *html.Node, s Selector) []*html.Node {
	result := make([]*html.Node, 0)

	recursion(n, s, &result)

	return result
}

func (q *Query) Find(n *html.Node) []*html.Node {
	result := make([]*html.Node, 0)

	if len(q.query) > 0 {
		input := []*html.Node{n}
		selector := q.query[0]

		result = selector.Filter(input)

		filters := q.query[1:]

		for i := range filters {
			if len(result) == 0 {
				break
			}
			result = filters[i].Filter(result)
		}
	}

	return result
}
