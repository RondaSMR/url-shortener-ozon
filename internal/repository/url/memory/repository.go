package memory

import (
	"context"
	apperor "url-shortener-ozon/internal/apperror"
	"url-shortener-ozon/internal/domain/entities"
)

func (r *Repository) CreateShortURL(_ context.Context, url entities.URLsStruct) error {
	r.pool[url.ShortURL] = url.OriginalURL
	return nil
}

func (r *Repository) GetShortURL(_ context.Context, url entities.InOutURL) (entities.InOutURL, error) {
	originalURL, ok := r.pool[url.URL]
	if !ok {
		return entities.InOutURL{}, apperor.ErrRepoNotFound
	}

	return entities.InOutURL{
		URL: originalURL,
	}, nil
}
