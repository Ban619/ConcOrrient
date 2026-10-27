package httpapi

import (
	"bytes"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"

	"concorrient/internal/workspace"
)

type Server struct {
	workspace *workspace.Service
}

func New(service *workspace.Service, webRoot fs.FS) http.Handler {
	server := &Server{workspace: service}
	router := chi.NewRouter()
	router.Get("/api/state", server.getState)
	router.Post("/api/ideas", server.createIdea)
	router.Post("/api/ideas/{id}/updates", server.addUpdate)
	router.Patch("/api/ideas/{id}", server.updateIdea)
	router.Post("/api/notices/read", server.readNotices)
	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		candidatePaths := []string{path}
		if path != "index.html" {
			candidatePaths = append(candidatePaths, "index.html")
		}

		for _, candidate := range candidatePaths {
			content, err := fs.ReadFile(webRoot, candidate)
			if err != nil {
				continue
			}

			info, err := fs.Stat(webRoot, candidate)
			if err != nil {
				continue
			}

			if ct := mime.TypeByExtension(filepath.Ext(candidate)); ct != "" {
				w.Header().Set("Content-Type", ct)
			}
			if candidate == "index.html" {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
			}

			http.ServeContent(w, r, candidate, info.ModTime(), bytes.NewReader(content))
			return
		}

		http.NotFound(w, r)
	}))
	return withHeaders(router)
}
