package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {
	// Placeholder for main function
	file, err := os.Open("ex3.html")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	doc, err := html.Parse(file)
	if err != nil {
		log.Fatalf("error parsing HTML: %v", err)
	}

	// bodyNode := bodyNode(doc)
	dfs(doc)

}

func dfs(body *html.Node) {
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "a" {
			printLinkNode(c)
		} else {
			dfs(c)
		}
	}
}

// func bodyNode(doc *html.Node) *html.Node {
// 	return doc.FirstChild.LastChild
// }

func printLinkNode(n *html.Node) {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			link := Link{
				Href: attr.Val,
				Text: extractText(n),
			}
			fmt.Printf("Link{\n  Href: %q,\n  Text: %q,\n}\n", link.Href, link.Text)
		}
	}
}

func extractText(n *html.Node) string {
	var parts []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		// include both text and comment nodes
		if n.Type == html.TextNode || n.Type == html.CommentNode {
			if s := strings.TrimSpace(n.Data); s != "" {
				parts = append(parts, s)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return strings.Join(parts, " ")
}
