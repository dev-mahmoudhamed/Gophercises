package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const (
	Limit      = 500
	MaxWorkers = 100 // max concurrent goroutines
)

func main() {
	http.HandleFunc("/", handleHome)
	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	totalStories, err := top500Stories()
	if err != nil {
		fmt.Println("error fetching top stories:", err)
		http.Error(w, "failed to fetch top stories", http.StatusInternalServerError)
		return
	}

	top := make(map[int]Story)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// create weighted semaphore
	sem := semaphore.NewWeighted(MaxWorkers)
	ctx := context.Background()

	for i := 0; i < Limit; i++ {
		idx := i
		storyID := totalStories[i]

		// acquire a worker slot
		if err := sem.Acquire(ctx, 1); err != nil {
			fmt.Println("semaphore acquire error:", err)
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer sem.Release(1) // release slot when done

			story, err := getStory(storyID)
			if err != nil {
				fmt.Println("error fetching story with id:", storyID, err)
				return
			}

			if story.Type == "story" && story.URL != "" {
				mu.Lock()
				top[idx] = story
				mu.Unlock()
			}
		}()
	}

	// wait for all goroutines to finish
	wg.Wait()

	duration := time.Since(start)

	type TemplateData struct {
		Stories      map[int]Story
		Duration     string
		TotalStories int
	}

	data := TemplateData{
		Stories:      top,
		Duration:     duration.String(),
		TotalStories: len(top),
	}

	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
