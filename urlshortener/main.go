package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	filePath := flag.String("file", "config.yaml", "Path to a YAML or JSON config file")
	flag.Parse()
	data, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	ext := filepath.Ext(*filePath)
	handler, err := fileHandler(ext, &data)

	if err != nil {
		log.Fatalf("error parsing file: %v", err)
	}

	if err != nil {
		log.Fatalf("error creating YAML handler: %v", err)
	}

	fmt.Println("🚀 Server running on http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}
