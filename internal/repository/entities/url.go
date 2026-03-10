package entities

import (
	"url-shortener-ozon/internal/domain/entities"
)

type RepoURL struct {
	URL string `json:"url"`
}

// RepoToEntity форматирует получаемую структуру из базы данных в обрабатываемую в коде
func (u RepoURL) RepoToEntity() entities.ResponseData {
	return entities.ResponseData{
		URL: u.URL,
	}
}
