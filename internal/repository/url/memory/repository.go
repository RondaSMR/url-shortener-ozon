package memory

import (
	"context"
	apperor "url-shortener-ozon/internal/apperror"
	"url-shortener-ozon/internal/domain/entities"
)

func (r *Repository) CreateShortPath(_ context.Context, url entities.URLsStruct) error {
	r.pool[url.ShortPath] = url.OriginalURL
	return nil
}

func (r *Repository) GetOriginalURLByShortPath(_ context.Context, url entities.RequestData) (entities.ResponseData, error) {
	originalURL, ok := r.pool[url.URL]
	if !ok {
		return entities.ResponseData{}, apperor.ErrRepoNotFound
	}

	return entities.ResponseData{
		URL: originalURL,
	}, nil
}
