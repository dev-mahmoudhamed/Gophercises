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

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <url>")
		fmt.Println("Please provide a URL as a command-line argument.")
		os.Exit(1)
	}
	url := os.Args[1]
	htmlCode, err := DownloadHTML(url)

	if err != nil {
		log.Fatalf("Failed to download: %v", err)
	}

	reader := strings.NewReader(htmlCode)
	doc, err := html.Parse(reader)

	if err != nil {
		log.Fatalf("error parsing HTML: %v", err)
	}

	links := make(map[string]bool)
	dfs(doc, links)

	sitemapXML, err := GenerateSitemap(links)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating sitemap: %v\n", err)
		return
	}

	fmt.Println(string(sitemapXML))

}
