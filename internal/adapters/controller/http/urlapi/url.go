package urlapi

import (
	"url-shortener-ozon/internal/domain/entities"
)

type InOutURL struct {
	URL string `json:"url"`
}

// AdapterHttpURLToEntity форматирует структуру получаемой информации в виде json запроса
func AdapterHttpURLToEntity(url InOutURL) (entities.InOutURL, error) {
	return entities.InOutURL{
		URL: url.URL,
	}, nil
}

// AdapterEntityToHttpURL преобразует структуру в формат json для отправки пользователю
func AdapterEntityToHttpURL(url entities.InOutURL) InOutURL {
	return InOutURL{
		URL: url.URL,
	}
}
