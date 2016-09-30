package parse

import (
    "golang.org/x/net/html"
    "strings"
)

type Entry struct {
    title string
    author string
    link string
}

func GetAttr(node *html.Node, name string) string {
    result := ""
    for _, a := range node.Attr {
        if a.Key == name {
            result = a.Val
            break
        }
    }

    return result
}

func GetClasses(node *html.Node) []string {
    var result []string = make([]string, 0)
    if node != nil {
        for _, a := range node.Attr {
            if a.Key == "class" {
                result = strings.Split(a.Val, " ")
                break
            }
        }
    }

    return result
}

func GetId(node *html.Node) string {
    return GetAttr(node, "id")
}

func GetHref(node *html.Node) string {
    return GetAttr(node, "href")
}