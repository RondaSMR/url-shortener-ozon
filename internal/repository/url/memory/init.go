package memory

type Repository struct {
	pool map[string]string
}

func NewMemoryRepository() *Repository {
	return &Repository{
		pool: make(map[string]string),
	}
}
