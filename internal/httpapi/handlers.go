package httpapi

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"concorrient/internal/workspace"
)

func (s *Server) getState(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.workspace.Snapshot(r.URL.Query().Get("actor")))
}

func (s *Server) createIdea(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Title        string   `json:"title"`
		Summary      string   `json:"summary"`
		AuthorID     string   `json:"authorId"`
		Mode         string   `json:"mode"`
		Collaborator string   `json:"collaborator"`
		Tags         []string `json:"tags"`
	}
	if !decodeJSON(w, r, &request) {
		return
	}
	input := workspace.CreateIdeaInput{
		Title: request.Title, Summary: request.Summary, AuthorID: request.AuthorID,
		Mode: request.Mode, Collaborator: request.Collaborator, Tags: request.Tags,
	}
	idea, err := s.workspace.CreateIdea(input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, idea)
}

func (s *Server) addUpdate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AuthorID string `json:"authorId"`
		Message  string `json:"message"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	idea, err := s.workspace.AddUpdate(chi.URLParam(r, "id"), input.AuthorID, input.Message)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, idea)
}

func (s *Server) updateIdea(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Status string `json:"status"`
		Actor  string `json:"actor"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	idea, err := s.workspace.SetStatus(chi.URLParam(r, "id"), input.Actor, input.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, idea)
}

func (s *Server) readNotices(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Actor string `json:"actor"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	if err := s.workspace.MarkNoticesRead(input.Actor); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeServiceError(w http.ResponseWriter, err error) {
	var validationError *workspace.ValidationError
	if errors.As(err, &validationError) {
		http.Error(w, validationError.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, workspace.ErrIdeaNotFound) {
		http.NotFound(w, nil)
		return
	}
	http.Error(w, "The workspace could not complete that request.", http.StatusInternalServerError)
}
