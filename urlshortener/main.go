package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {

	// pathsToUrls := map[string]string{
	// 	"/google": "https://google.com",
	// 	"/yahoo":  "https://yahoo.com",
	// }
	// mapHandler := MapHandler(pathsToUrls)

	yamlData := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

	yamlData = strings.TrimSpace(yamlData)

	yamlHandler, err := YAMLHandler([]byte(yamlData))
	if err != nil {
		log.Fatalf("error creating YAML handler: %v", err)
	}
	fmt.Println("🚀 Server running on http://localhost:8080")
	http.ListenAndServe(":8080", yamlHandler)
}
