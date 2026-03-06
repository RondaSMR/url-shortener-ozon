package usecase

import (
	"context"
	"errors"
	"fmt"
	apperor "url-shortener-ozon/internal/apperror"
	"url-shortener-ozon/internal/controller/http/v1/urlshortener"
	"url-shortener-ozon/internal/domain/entities"
	"url-shortener-ozon/pkg/utils"
)

var _ urlshortener.Usecase = new(Usecase)

type Repository interface {
	CreateShortURL(ctx context.Context, url entities.URLsStruct) error
	GetShortURL(ctx context.Context, url entities.InOutURL) (entities.InOutURL, error)
}

type Usecase struct {
	urlRepository Repository
}

func NewUseCase(
	urlRepository Repository,
) *Usecase {
	return &Usecase{urlRepository: urlRepository}
}

func (u Usecase) CreateShortURL(ctx context.Context, url *entities.InOutURL) (entities.InOutURL, error) {
	// Проверка валидации URL адреса
	if err := url.Validate(); err != nil {
		return *url, fmt.Errorf("validation failed: %w", err)
	}

	// Пробуем с разной солью при коллизии
	for salt := 0; ; salt++ {
		// Генерация короткой ссылки
		shortURL := utils.GenerateShortURL(*url, salt)

		// Если по сгенерированной ссылке нашлась другая (длинная) и нет ошибок
		if longURL, err := u.urlRepository.GetShortURL(ctx, shortURL); err == nil {
			if longURL.URL == url.URL {
				// Уже существует — возвращаем
				return shortURL, nil
			}
			// Коллизия — пробуем другую соль
			continue
			// Если вернулась ошибка о ненайденном значении по ключу сгенерированной ссылки
		} else if errors.Is(err, apperor.ErrRepoNotFound) {
			// Ссылка уникальна - возвращаем
			return shortURL, u.urlRepository.CreateShortURL(ctx, entities.URLsStruct{
				OriginalURL: url.URL,
				ShortURL:    shortURL.URL,
			})
		} else {
			// Была найдена другая ошибка
			return *url, err
		}
	}
}

func (u Usecase) GetShortURL(ctx context.Context, url *entities.InOutURL) (entities.InOutURL, error) {
	return u.urlRepository.GetShortURL(ctx, *url)
}
