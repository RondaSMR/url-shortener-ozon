package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	apperor "url-shortener-ozon/internal/apperror"
	"url-shortener-ozon/internal/domain/entities"
	repoEntities "url-shortener-ozon/internal/repository/entities"
)

func (r Repository) CreateShortPath(ctx context.Context, url entities.URLsStruct) error {
	_, err := r.pool.Exec(ctx, `insert into urls(original_url, short_path) values ($1, $2)`, url.OriginalURL, url.ShortPath)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetOriginalURLByShortPath(ctx context.Context, url entities.RequestData) (entities.ResponseData, error) {
	row, err := r.pool.Query(ctx, `select original_url as url from urls where short_path = $1`, url.URL)
	if err != nil {
		return entities.ResponseData{}, fmt.Errorf("executing query error: %w", err)
	}
	defer row.Close()

	urlRepo, err := pgx.CollectOneRow(row, pgx.RowToStructByName[repoEntities.RepoURL])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.ResponseData{}, apperor.ErrRepoNotFound
		}
		return entities.ResponseData{}, fmt.Errorf("selecting url from database error: %w", err)
	}

	return urlRepo.RepoToEntity(), nil
}
