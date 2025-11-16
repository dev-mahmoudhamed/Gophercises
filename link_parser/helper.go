package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type sitemapURL struct {
	Loc string `xml:"loc"`
}

type urlset struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

func dfs(body *html.Node, links map[string]bool) {
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "a" {
			path := strings.Trim(linkPath(c), "/")
			links[path] = false
			// *links = append(*links, path)
		} else {
			dfs(c, links)
		}
	}
}

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

func linkPath(n *html.Node) string {
	var path string
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			path = attr.Val
		}
	}
	return path
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

func GenerateSitemap(links map[string]bool) ([]byte, error) {

	// 1. Create the root <urlset> structure
	sitemap := &urlset{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]sitemapURL, 0, len(links)),
	}

	// 2. Populate the <url> entries from the string slice
	for link := range links {
		//sitemap.URLs[i] = sitemapURL{Loc: link}
		sitemap.URLs = append(sitemap.URLs, sitemapURL{Loc: link})
	}

	// 3. Marshal the data into XML format.
	// We use MarshalIndent for pretty, human-readable output.
	xmlData, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		return nil, err
	}

	// 4. Prepend the standard XML header
	// The xml.Marshal functions do not add this by default.
	output := []byte(xml.Header)
	output = append(output, '\n')
	output = append(output, xmlData...)

	return output, nil
}

func DownloadHTML(url string) (string, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Return a formatted error
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
