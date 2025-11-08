package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

var story Adventure
var tmpl *template.Template

func main() {
	var err error
	story, err = ParseJson("gopher.json")
	if err != nil {
		log.Fatal(err)
	}

	tmpl = template.Must(template.ParseFiles("story.html"))

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, "/intro", http.StatusFound)
	// })

	// Handle story chapters
	http.HandleFunc("/", handleChapter)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleChapter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	if path == "" {
		http.Redirect(w, r, "/intro", http.StatusFound)
		return
	}

	if chapter, ok := story[path]; ok {
		err := tmpl.Execute(w, chapter)
		if err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Chapter not found", http.StatusNotFound)
	}
}
