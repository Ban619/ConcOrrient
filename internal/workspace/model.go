package workspace

import "time"

type Member struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Initials string `json:"initials"`
	Color    string `json:"color"`
}

type Update struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"authorId"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type Idea struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	Status       string    `json:"status"`
	Mode         string    `json:"mode"`
	AuthorID     string    `json:"authorId"`
	Collaborator string    `json:"collaborator,omitempty"`
	Tags         []string  `json:"tags"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Updates      []Update  `json:"updates"`
}

type Notice struct {
	ID        string    `json:"id"`
	MemberID  string    `json:"memberId"`
	IdeaID    string    `json:"ideaId"`
	Text      string    `json:"text"`
	Kind      string    `json:"kind"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"createdAt"`
}

type State struct {
	Members []Member `json:"members"`
	Ideas   []Idea   `json:"ideas"`
	Notices []Notice `json:"notices"`
}

type View struct {
	Members []Member `json:"members"`
	Ideas   []Idea   `json:"ideas"`
	Notices []Notice `json:"notices"`
	Actor   string   `json:"actor"`
}

type CreateIdeaInput struct {
	Title        string
	Summary      string
	AuthorID     string
	Mode         string
	Collaborator string
	Tags         []string
}
