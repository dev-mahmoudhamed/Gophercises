package main

import (
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

func YAMLHandler(yamlBytes []byte) (yamlHandler http.HandlerFunc, err error) {
	parsedYaml, err := parseYAML(yamlBytes)
	if err != nil {
		return
	}
	pathMap := buildMap(parsedYaml)
	yamlHandler = MapHandler(pathMap)
	return
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

type pathToURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
