package main

import (
	"encoding/json"
	"os"
)

type Arc struct {
	Title   string
	Story   []string
	Options []struct {
		Text string
		Arc  string
	}
}
type Adventure map[string]Arc

func ParseJson(path string) (Adventure, error) {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var story Adventure
	if err := json.Unmarshal(jsonData, &story); err != nil {
		return nil, err
	}
	return story, nil
}
