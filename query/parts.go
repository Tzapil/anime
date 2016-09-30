package query

import (
    "golang.org/x/net/html"
    "github.com/tzapil/anime/parse"
    "strings"
)

// Interface
type Selector interface {
    Check(n *html.Node) bool
    Filter(n []*html.Node) []*html.Node
}

// Universal
type universal struct {}

func (u *universal) Check(n *html.Node) bool {
    return true
}

func (u *universal) Filter(n []*html.Node) []*html.Node {
    return Filter(n, u)
}

func NewUniversal() *universal {
    return &universal{}
}

// Child
type child struct {
    universal
    s Selector
}

func (u *child) Check(n *html.Node) bool {
    return u.s.Check(n)
}

func (u *child) Filter(n []*html.Node) []*html.Node {
    result := make([]*html.Node, 0)
    for _, v := range n {
        child := v.FirstChild
        for child != nil {
            if u.Check(child) {
                result = append(result, child)
            }
            child = child.NextSibling
        }
    }

    return result
}

func NewChild(s Selector) *child {
    return &child{*NewUniversal(), s}
}

// Next
type next struct {
    universal
}

func (u *next) Filter(n []*html.Node) []*html.Node {
    result := make([]*html.Node, 0)
    for _, v := range n {
        nxt := v.NextSibling
        if nxt != nil {
            result = append(result, nxt)
        }
    }

    return result
}

func NewNext() *next {
    return &next{*NewUniversal()}
}

// Contains
type contains struct {
    universal
    val string
}

func (u *contains) Check(n *html.Node) bool {
    return n != nil && strings.Contains(Text(n), u.val)
}

func (u *contains) Filter(n []*html.Node) []*html.Node {
    return Filter(n, u)
}

func NewContains(val string) *contains{
    return &contains{*NewUniversal(), val}
}

// Tag Selector
type tag struct {
    universal
    name string
}

func (u *tag) Check(n *html.Node) bool {
    return n != nil && u.name == n.Data
}

func NewTag(name string) *tag{
    return &tag{*NewUniversal(), name}
}

// Class Selector
type class struct {
    universal
    name string
}

func (u *class) Check(n *html.Node) bool {
    return n != nil && Contains(u.name, parse.GetClasses(n))
}

func NewClass(name string) *class{
    return &class{*NewUniversal(), name}
}