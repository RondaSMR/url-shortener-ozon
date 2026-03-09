package entities

import (
	"fmt"
	"net/url"
)

type InOutURL struct {
	URL string
}

type URLsStruct struct {
	OriginalURL string
	ShortURL    string
}

// Validate проверяет поступивший URL
func (u InOutURL) Validate() error {
	parsedURL, err := url.ParseRequestURI(u.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must be http:// or https://")
	}

	return nil
}
