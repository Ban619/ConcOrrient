package workspace

import (
	"reflect"
	"strings"
	"testing"
)

type memoryStore struct {
	state State
}

func (s *memoryStore) Load() (State, error) {
	return s.state, nil
}

func (s *memoryStore) Save(state State) error {
	s.state = state
	return nil
}

func TestSetStatusNotifiesRelevantMembersOnlyWhenChanged(t *testing.T) {
	tests := []struct {
		name       string
		ideaID     string
		actor      string
		status     string
		recipients []string
		noticeKind string
	}{
		{
			name:   "open thread notifies the workspace",
			ideaID: "thread-momentum", actor: "alex", status: "exploring",
			recipients: []string{"mina", "jun", "noor"}, noticeKind: "open",
		},
		{
			name:   "direct thread notifies its collaborator",
			ideaID: "thread-notes", actor: "mina", status: "exploring",
			recipients: []string{"noor"}, noticeKind: "direct",
		},
		{
			name:   "unchanged status does not notify",
			ideaID: "thread-handoffs", actor: "jun", status: "exploring",
			recipients: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := &memoryStore{state: seedState()}
			service, err := NewService(store)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := service.SetStatus(test.ideaID, test.actor, test.status); err != nil {
				t.Fatal(err)
			}

			recipients := make([]string, 0)
			for _, notice := range store.state.Notices {
				recipients = append(recipients, notice.MemberID)
				if test.noticeKind != "" && notice.Kind != test.noticeKind {
					t.Errorf("notice kind = %q, want %q", notice.Kind, test.noticeKind)
				}
				if !strings.Contains(notice.Text, "changed status to "+test.status) {
					t.Errorf("notice text %q does not describe the status change", notice.Text)
				}
			}
			if !reflect.DeepEqual(recipients, test.recipients) {
				t.Errorf("notice recipients = %v, want %v", recipients, test.recipients)
			}
		})
	}
}
