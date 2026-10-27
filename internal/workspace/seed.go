package workspace

import "time"

func seedState() State {
	now := time.Now().UTC()
	return State{
		Members: []Member{
			{ID: "mina", Name: "Mina Patel", Initials: "MP", Color: "lilac"},
			{ID: "alex", Name: "Alex Rivera", Initials: "AR", Color: "coral"},
			{ID: "jun", Name: "Jun Park", Initials: "JP", Color: "mint"},
			{ID: "noor", Name: "Noor Hassan", Initials: "NH", Color: "blue"},
		},
		Ideas: []Idea{
			{
				ID: "thread-handoffs", Title: "What makes a handoff feel effortless?",
				Summary: "A good idea shouldn’t lose its shape when someone signs off. What context helps the next person keep going?",
				Status:  "exploring", Mode: "open", AuthorID: "jun", Tags: []string{"handoffs", "workflow"},
				CreatedAt: now.Add(-7 * time.Hour), UpdatedAt: now.Add(-42 * time.Minute),
				Updates: []Update{
					{ID: "u-h1", AuthorID: "jun", Message: "Starting with a question: what is the smallest useful handoff note?", CreatedAt: now.Add(-7 * time.Hour)},
					{ID: "u-h2", AuthorID: "alex", Message: "A decision and the reason behind it would help me pick up where someone left off.", CreatedAt: now.Add(-42 * time.Minute)},
				},
			},
			{
				ID: "thread-notes", Title: "Give unfinished thoughts somewhere to land",
				Summary: "Let people leave a useful seed without needing to turn it into a polished proposal first.",
				Status:  "open", Mode: "direct", AuthorID: "mina", Collaborator: "noor", Tags: []string{"ideas", "culture"},
				CreatedAt: now.Add(-5 * time.Hour), UpdatedAt: now.Add(-2 * time.Hour),
				Updates: []Update{
					{ID: "u-n1", AuthorID: "mina", Message: "Could a thread start as a rough observation and grow as people add context?", CreatedAt: now.Add(-5 * time.Hour)},
					{ID: "u-n2", AuthorID: "noor", Message: "Yes. I’d also want the next useful question to be visible at a glance.", CreatedAt: now.Add(-2 * time.Hour)},
				},
			},
			{
				ID: "thread-momentum", Title: "Keep momentum across different schedules",
				Summary: "A gentle signal when a teammate has added something worth looking at, without expecting an instant reply.",
				Status:  "open", Mode: "open", AuthorID: "alex", Tags: []string{"async", "signals"},
				CreatedAt: now.Add(-24 * time.Hour), UpdatedAt: now.Add(-11 * time.Hour),
				Updates: []Update{
					{ID: "u-m1", AuthorID: "alex", Message: "Maybe activity can be a quiet invitation instead of another urgent ping.", CreatedAt: now.Add(-24 * time.Hour)},
					{ID: "u-m2", AuthorID: "mina", Message: "A short recap in the notification would give me enough context to decide when to jump in.", CreatedAt: now.Add(-11 * time.Hour)},
				},
			},
		},
		Notices: []Notice{},
	}
}
