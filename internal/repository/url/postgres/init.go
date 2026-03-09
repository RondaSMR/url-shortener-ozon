package postgres

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"url-shortener-ozon/pkg/config"
	"url-shortener-ozon/pkg/connectors/pgconnector"
)

type Repository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository инициализирует Postgres-хранилище
func NewPostgresRepository(config config.PgStorage) (*Repository, error) {

	// Добавление конфигурации PostgreSQL
	connectorConfig, err := pgconnector.CreateConfig(&pgconnector.ConnectionConfig{
		Host:     config.Host,
		Port:     fmt.Sprint(config.Port),
		User:     config.User,
		Password: config.Pass,
		DbName:   config.DB,
		SslMode:  "disable",
	},
		nil)

	// Инициализация PostgreSQL
	pgStorage, err := pgconnector.NewPgConnector(
		connectorConfig,
		10*time.Second,
		10*time.Second,
	)

	if err != nil {
		return nil, fmt.Errorf("initialize pg storage: %w", err)
	}

	return &Repository{
		pool: pgStorage.GetPool(),
	}, nil
}

// Close нужен для закрытия pool внутри функции NewAppPostgres
func (r *Repository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}
