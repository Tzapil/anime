package query

import (
    "golang.org/x/net/html"
    "bytes"
)

// Helpers
func Filter(n []*html.Node, s Selector) []*html.Node {
    result := make([]*html.Node, 0)
    for i := range n {
        if s.Check(n[i]) {
            result = append(result, n[i])
        }
    }

    return result
}

func Contains(n string, arr []string) bool {
    result := false
    for _, a := range arr {
        if a == n {
            result = true
            break
        }
    }

    return result
}

func Text(n *html.Node) string {
    result := ""

    if n != nil {
        wr := new(bytes.Buffer)
        html.Render(wr, n)
        result = wr.String()
    }

    return result
}