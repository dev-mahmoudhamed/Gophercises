package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {

	handler, err := appHandler()

	if err != nil {
		log.Fatalf("error parsing file: %v", err)
	}

	fmt.Println("🚀 Server running on http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}

func appHandler() (http.HandlerFunc, error) {
	fileFlag := flag.String("file", "", "Path to a YAML or JSON config file")
	redisAddr := flag.String("db", "", "Redis server address")
	flag.Parse()

	switch {
	case *fileFlag != "":
		return fileHandler(*fileFlag)
	case *redisAddr != "":
		return redisHandler(*redisAddr)
	default:
		return nil, fmt.Errorf("no configuration source provided (use -file or -db)")
	}

}
