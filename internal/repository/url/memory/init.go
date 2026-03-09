package memory

type Repository struct {
	pool map[string]string
}

// NewMemoryRepository инициализирует map-хранилище
func NewMemoryRepository() *Repository {
	return &Repository{
		pool: make(map[string]string),
	}
}
