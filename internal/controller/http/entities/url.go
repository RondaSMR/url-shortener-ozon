package entities

import (
	"url-shortener-ozon/internal/domain/entities"
)

// RequestDTOData - DTO для HTTP слоя
type RequestDTOData struct {
	URL string `json:"url"`
}

type ResponseDTOData struct {
	URL string `json:"url"`
}

// ToEntity преобразует DTO в доменную сущность
func (u RequestDTOData) ToEntity() entities.RequestData {
	return entities.RequestData{
		URL: u.URL,
	}
}

// FromEntity преобразует доменную сущность в DTO
func FromEntity(url entities.ResponseData) ResponseDTOData {
	return ResponseDTOData{
		URL: url.URL,
	}
}
