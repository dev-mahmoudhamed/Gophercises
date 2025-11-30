package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Story struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

func top500Stories() (map[int]int, error) {
	stories := make(map[int]int)

	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, fmt.Errorf("error fetching top stories: %w", err)
	}
	defer resp.Body.Close()

	var result []int
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	for i, v := range result {
		stories[i] = v
	}
	return stories, nil
}

func getStory(id int) (Story, error) {
	story := Story{}

	resp, err := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id))
	if err != nil {
		return story, fmt.Errorf("error fetching story: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
		return story, err
	}
	return story, nil
}
