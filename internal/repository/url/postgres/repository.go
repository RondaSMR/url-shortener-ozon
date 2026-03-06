package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"url-shortener-ozon/internal/adapters/repository"
	apperor "url-shortener-ozon/internal/apperror"
	"url-shortener-ozon/internal/domain/entities"
)

func (r Repository) CreateShortURL(ctx context.Context, url entities.URLsStruct) error {
	_, err := r.pool.Exec(ctx, `insert into urls(original_url, short_url) values ($1, $2)`, url.OriginalURL, url.ShortURL)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetShortURL(ctx context.Context, url entities.InOutURL) (entities.InOutURL, error) {
	row, err := r.pool.Query(ctx, `select original_url as url from urls where short_url = $1`, url.URL)
	if err != nil {
		return entities.InOutURL{}, fmt.Errorf("executing query error: %w", err)
	}
	defer row.Close()

	urlRepo, err := pgx.CollectOneRow(row, pgx.RowToStructByName[repository.InOutURL])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.InOutURL{}, apperor.ErrRepoNotFound
		}
		return entities.InOutURL{}, fmt.Errorf("selecting url from database error: %w", err)
	}

	return repository.AdapterRepoTaskToEntity(urlRepo), nil
}
