package workspace

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"sync"
	"time"
)

var ErrIdeaNotFound = errors.New("thread not found")

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type Service struct {
	mu    sync.RWMutex
	state State
	store Store
}

func NewService(store Store) (*Service, error) {
	state, err := store.Load()
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		state = seedState()
		if err := store.Save(state); err != nil {
			return nil, err
		}
	}
	return &Service{state: state, store: store}, nil
}

func (s *Service) Snapshot(actor string) View {
	s.mu.RLock()
	defer s.mu.RUnlock()

	view := View{
		Members: append([]Member(nil), s.state.Members...),
		Ideas:   make([]Idea, len(s.state.Ideas)),
		Notices: make([]Notice, 0),
		Actor:   actor,
	}
	if !s.isMember(actor) && len(s.state.Members) > 0 {
		view.Actor = s.state.Members[0].ID
	}
	for i, idea := range s.state.Ideas {
		view.Ideas[i] = cloneIdea(idea)
	}
	for _, notice := range s.state.Notices {
		if notice.MemberID == view.Actor {
			view.Notices = append(view.Notices, notice)
		}
	}
	return view
}

func (s *Service) CreateIdea(input CreateIdeaInput) (Idea, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Summary = strings.TrimSpace(input.Summary)
	if input.Title == "" || input.Summary == "" || len(input.Title) > 120 || len(input.Summary) > 700 {
		return Idea{}, &ValidationError{Message: "Add a title and a short description (up to 700 characters)."}
	}
	if input.Mode != "direct" && input.Mode != "open" {
		return Idea{}, &ValidationError{Message: "Choose direct or open collaboration."}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.isMember(input.AuthorID) {
		return Idea{}, &ValidationError{Message: "Choose a workspace member."}
	}
	if input.Mode == "direct" && (!s.isMember(input.Collaborator) || input.Collaborator == input.AuthorID) {
		return Idea{}, &ValidationError{Message: "Choose a teammate to invite."}
	}

	now := time.Now().UTC()
	idea := Idea{
		ID: newID("idea"), Title: input.Title, Summary: input.Summary, Status: "open", Mode: input.Mode,
		AuthorID: input.AuthorID, Collaborator: input.Collaborator, Tags: cleanTags(input.Tags),
		CreatedAt: now, UpdatedAt: now,
		Updates: []Update{{ID: newID("update"), AuthorID: input.AuthorID, Message: input.Summary, CreatedAt: now}},
	}
	s.state.Ideas = append([]Idea{idea}, s.state.Ideas...)
	s.notifyForIdea(idea, input.AuthorID, input.Mode == "direct", "started a thread")
	if err := s.saveLocked(); err != nil {
		return Idea{}, err
	}
	return cloneIdea(idea), nil
}

func (s *Service) AddUpdate(ideaID, authorID, message string) (Idea, error) {
	message = strings.TrimSpace(message)
	if message == "" || len(message) > 1000 {
		return Idea{}, &ValidationError{Message: "Write an update (up to 1,000 characters)."}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.isMember(authorID) {
		return Idea{}, &ValidationError{Message: "Choose a workspace member."}
	}
	index := s.ideaIndex(ideaID)
	if index < 0 {
		return Idea{}, ErrIdeaNotFound
	}

	now := time.Now().UTC()
	idea := &s.state.Ideas[index]
	idea.Updates = append(idea.Updates, Update{ID: newID("update"), AuthorID: authorID, Message: message, CreatedAt: now})
	idea.UpdatedAt = now
	s.notifyForIdea(*idea, authorID, idea.Mode == "direct", "added to a thread")
	if err := s.saveLocked(); err != nil {
		return Idea{}, err
	}
	return cloneIdea(*idea), nil
}

func (s *Service) SetStatus(ideaID, actor, status string) (Idea, error) {
	if status != "open" && status != "exploring" && status != "resolved" {
		return Idea{}, &ValidationError{Message: "Choose open, exploring, or resolved."}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	index := s.ideaIndex(ideaID)
	if index < 0 {
		return Idea{}, ErrIdeaNotFound
	}
	if !s.isMember(actor) {
		return Idea{}, &ValidationError{Message: "Choose a workspace member."}
	}

	idea := &s.state.Ideas[index]
	statusChanged := idea.Status != status
	idea.Status = status
	idea.UpdatedAt = time.Now().UTC()
	if statusChanged {
		s.notifyForIdea(*idea, actor, idea.Mode == "direct", "changed status to "+status)
	}
	if err := s.saveLocked(); err != nil {
		return Idea{}, err
	}
	return cloneIdea(*idea), nil
}

func (s *Service) MarkNoticesRead(actor string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.isMember(actor) {
		return &ValidationError{Message: "Choose a workspace member."}
	}
	for i := range s.state.Notices {
		if s.state.Notices[i].MemberID == actor {
			s.state.Notices[i].Read = true
		}
	}
	return s.saveLocked()
}

func (s *Service) saveLocked() error {
	return s.store.Save(s.state)
}

func (s *Service) notifyForIdea(idea Idea, actor string, directOnly bool, action string) {
	var author string
	for _, member := range s.state.Members {
		if member.ID == actor {
			author = member.Name
			break
		}
	}
	for _, member := range s.state.Members {
		if member.ID == actor || (directOnly && member.ID != idea.AuthorID && member.ID != idea.Collaborator) {
			continue
		}
		kind := "open"
		if directOnly {
			kind = "direct"
		}
		s.state.Notices = append(s.state.Notices, Notice{
			ID: newID("notice"), MemberID: member.ID, IdeaID: idea.ID,
			Text: author + " " + action + ": " + idea.Title, Kind: kind, CreatedAt: time.Now().UTC(),
		})
	}
}

func (s *Service) isMember(id string) bool {
	if id == "" {
		return false
	}
	for _, member := range s.state.Members {
		if member.ID == id {
			return true
		}
	}
	return false
}

func (s *Service) ideaIndex(id string) int {
	for i := range s.state.Ideas {
		if s.state.Ideas[i].ID == id {
			return i
		}
	}
	return -1
}

func cleanTags(tags []string) []string {
	cleaned := make([]string, 0, 3)
	seen := make(map[string]bool)
	for _, tag := range tags {
		tag = strings.TrimSpace(strings.TrimPrefix(tag, "#"))
		if tag == "" || len(tag) > 24 || seen[strings.ToLower(tag)] || len(cleaned) == 3 {
			continue
		}
		seen[strings.ToLower(tag)] = true
		cleaned = append(cleaned, tag)
	}
	return cleaned
}

func newID(prefix string) string {
	var value [8]byte
	if _, err := rand.Read(value[:]); err != nil {
		return fmt.Sprintf("%s-%s", prefix, time.Now().UTC().Format("20060102150405"))
	}
	return prefix + "-" + hex.EncodeToString(value[:])
}

func cloneIdea(idea Idea) Idea {
	idea.Tags = append([]string(nil), idea.Tags...)
	idea.Updates = append([]Update(nil), idea.Updates...)
	return idea
}
