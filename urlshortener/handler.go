package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

func MapHandler(pathsToURLs map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dest, ok := pathsToURLs[r.URL.Path]; ok && dest != "" {
			http.Redirect(w, r, dest, http.StatusPermanentRedirect)
			return
		}
		http.NotFound(w, r)
	}
}

func fileHandler(ext string, data *[]byte) (http.HandlerFunc, error) {
	var parsed []pathToURL
	var err error

	switch ext {
	case ".yaml", ".yml":
		parsed, err = parseYAML(*data)
	case ".json":
		parsed, err = parseJSON(*data)
	default:
		log.Fatalf("unsupported file type: %s", ext)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}
	pathMap := buildMap(parsed)
	return MapHandler(pathMap), err
}

func buildMap(pathsToURLs []pathToURL) (builtMap map[string]string) {
	builtMap = make(map[string]string)
	for _, ptu := range pathsToURLs {
		builtMap[ptu.Path] = ptu.URL
	}
	return
}

func parseYAML(yamlData []byte) (pathsToURLs []pathToURL, err error) {
	err = yaml.Unmarshal(yamlData, &pathsToURLs)
	return
}

func parseJSON(jsonData []byte) (pathsToURLs []pathToURL, err error) {
	err = json.Unmarshal(jsonData, &pathsToURLs)
	return
}

type pathToURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url"  json:"url"`
}
