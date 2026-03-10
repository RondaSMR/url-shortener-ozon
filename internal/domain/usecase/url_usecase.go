package usecase

import (
	"context"
	"errors"
	"fmt"
	apperor "url-shortener-ozon/internal/apperror"
	"url-shortener-ozon/internal/controller/http/v1/url_shortener"
	"url-shortener-ozon/internal/domain/entities"
	"url-shortener-ozon/pkg/utils"
)

var _ url_shortener.Usecase = new(Usecase)

type Repository interface {
	CreateShortPath(ctx context.Context, url entities.URLsStruct) error
	GetOriginalURLByShortPath(ctx context.Context, url entities.RequestData) (entities.ResponseData, error)
}

type Usecase struct {
	urlRepository Repository
}

func NewUseCase(
	urlRepository Repository,
) *Usecase {
	return &Usecase{urlRepository: urlRepository}
}

// CreateShortPath является бизнес-логикой процесса добавления сокращенной URL ссылки
func (u Usecase) CreateShortPath(ctx context.Context, requestData *entities.RequestData) (entities.ResponseData, error) {
	// Проверка валидации URL адреса
	if err := requestData.Validate(); err != nil {
		return entities.ResponseData{}, fmt.Errorf("validation failed: %w", err)
	}

	// Пробуем с разной солью при коллизии
	for salt := 0; ; salt++ {
		// Генерация короткой ссылки
		shortPath := utils.GenerateShortPath(requestData.URL, salt)

		// Если по сгенерированной ссылке нашлась другая (длинная) и нет ошибок
		if longURL, err := u.urlRepository.GetOriginalURLByShortPath(ctx, entities.RequestData{URL: shortPath}); err == nil {
			if longURL.URL == requestData.URL {
				// Уже существует — возвращаем
				return entities.ResponseData{URL: shortPath}, nil
			}
			// Коллизия — пробуем другую соль
			continue
			// Если вернулась ошибка о ненайденном значении по ключу сгенерированной ссылки
		} else if errors.Is(err, apperor.ErrRepoNotFound) {
			// Ссылка уникальна - возвращаем
			return entities.ResponseData{URL: shortPath}, u.urlRepository.CreateShortPath(ctx, entities.URLsStruct{
				OriginalURL: requestData.URL,
				ShortPath:   shortPath,
			})
		} else {
			// Была найдена другая ошибка
			return entities.ResponseData{}, err
		}
	}
}

// GetOriginalURLByShortPath является бизнес-логикой процесса перехода по сокращенной URL ссылки
func (u Usecase) GetOriginalURLByShortPath(ctx context.Context, url *entities.RequestData) (entities.ResponseData, error) {
	return u.urlRepository.GetOriginalURLByShortPath(ctx, *url)
}
