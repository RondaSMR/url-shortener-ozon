package repository

import (
	"database/sql"
	"url-shortener-ozon/internal/domain/entities"
)

type InOutURL struct {
	URL sql.NullString `json:"url"`
}

// AdapterRepoTaskToEntity форматирует получаемую структуру из базы данных в обрабатываемую в коде
func AdapterRepoTaskToEntity(url InOutURL) entities.InOutURL {
	return entities.InOutURL{
		URL: url.URL.String,
	}
}
