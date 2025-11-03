package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/redis/go-redis/v9"
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

func fileHandler(flag string) (http.HandlerFunc, error) {

	data, oserr := os.ReadFile(flag)
	if oserr != nil {
		log.Fatalf("error reading file: %v", oserr)
	}

	ext := filepath.Ext(flag)

	var parsed []pathToURL
	var err error

	switch ext {
	case ".yaml", ".yml":
		parsed, err = parseYAML(data)
	case ".json":
		parsed, err = parseJSON(data)
	default:
		log.Fatalf("unsupported file type: %s", ext)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}
	pathMap := buildMap(parsed)
	return MapHandler(pathMap), err
}

func redisHandler(addr string) (http.HandlerFunc, error) {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0, // default
	})

	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	pathMap := make(map[string]string)
	for _, key := range keys {
		url, err := rdb.Get(ctx, key).Result()
		if err == nil {
			pathMap[key] = url
		}
	}

	return MapHandler(pathMap), nil
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
