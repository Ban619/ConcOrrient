package workspace

type Store interface {
	Load() (State, error)
	Save(State) error
}
