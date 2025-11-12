package main

import (
	"flag"
	"log"
	"os"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {

	var filePath string
	flag.StringVar(&filePath, "file", "ex3.html", "path to HTML file")
	flag.Parse()

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	doc, err := html.Parse(file)
	if err != nil {
		log.Fatalf("error parsing HTML: %v", err)
	}
	dfs(doc)
}
