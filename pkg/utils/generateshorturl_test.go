package utils

import (
	"testing"
	"url-shortener-ozon/internal/domain/entities"
)

func TestUtils_GenerateShortURL(t *testing.T) {
	originalURL := entities.InOutURL{URL: "https://ozon.ru"}

	shortURL := GenerateShortURL(originalURL, 0)

	if shortURL.URL != GenerateShortURL(originalURL, 0).URL {
		t.Error("Generated different short URLs for the same original URL")
	}
}
