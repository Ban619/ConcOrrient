package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"concorrient/internal/httpapi"
	"concorrient/internal/workspace"
)

//go:embed web/dist/* web/dist/_astro/*
var webFiles embed.FS

func main() {
	dataPath := os.Getenv("CONCORRIENT_DATA")
	if dataPath == "" {
		dataPath = "data.json"
	}

	service, err := workspace.NewService(workspace.NewJSONStore(dataPath))
	if err != nil {
		log.Fatal(err)
	}

	webRoot, err := fs.Sub(webFiles, "web/dist")
	if err != nil {
		log.Fatal(err)
	}

	addr := os.Getenv("CONCORRIENT_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("ConcOrrient listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, httpapi.New(service, webRoot)))
}
