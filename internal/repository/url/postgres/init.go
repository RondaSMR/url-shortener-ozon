package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"url-shortener-ozon/pkg/connectors/pgconnector"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pgConnector *pgconnector.Connector) *Repository {
	return &Repository{
		pool: pgConnector.GetPool(),
	}
}
