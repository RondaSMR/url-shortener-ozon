package postgres_test

import (
	"context"
	"testing"
	apperor "url-shortener-ozon/internal/apperror"

	"url-shortener-ozon/internal/domain/entities"
	"url-shortener-ozon/internal/domain/usecase/postgres"
)

type MockRepository struct {
	CreateFunc func(ctx context.Context, url entities.URLsStruct) error
	GetFunc    func(ctx context.Context, url entities.InOutURL) (entities.InOutURL, error)
}

func (m *MockRepository) CreateShortURL(ctx context.Context, url entities.URLsStruct) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, url)
	}
	return nil
}

func (m *MockRepository) GetShortURL(ctx context.Context, url entities.InOutURL) (entities.InOutURL, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, url)
	}
	return entities.InOutURL{}, nil
}

func TestPostgresUsecase_CreateAndGet(t *testing.T) {

	mockRepo := &MockRepository{
		CreateFunc: func(ctx context.Context, url entities.URLsStruct) error {
			if url.OriginalURL != "https://ozon.ru" {
				t.Errorf("Expected OriginalURL 'https://ozon.ru', got %s", url.OriginalURL)
			}
			if url.ShortURL == "" {
				t.Error("ShortURL is empty")
			}
			return nil
		},
		GetFunc: func(ctx context.Context, url entities.InOutURL) (entities.InOutURL, error) {
			return entities.InOutURL{URL: "https://ozon.ru"}, nil
		},
	}

	usecase := postgres.NewUseCase(mockRepo)
	ctx := context.Background()

	postURL := entities.InOutURL{URL: "https://ozon.ru"}
	shortURL, err := usecase.CreateShortURL(ctx, &postURL)
	if err != nil {
		t.Fatalf("CreateShortURL failed: %v", err)
	}
	if shortURL.URL == "" {
		t.Fatal("CreateShortURL returned empty URL")
	}

	getURL := entities.InOutURL{URL: shortURL.URL}
	originalURL, err := usecase.GetShortURL(ctx, &getURL)
	if err != nil {
		t.Fatalf("GetShortURL failed: %v", err)
	}
	if originalURL.URL != "https://ozon.ru" {
		t.Fatalf("Expected 'https://ozon.ru', got %s", originalURL.URL)
	}
}

func TestPostgresUsecase_NotFound(t *testing.T) {

	mockRepo := &MockRepository{
		GetFunc: func(ctx context.Context, url entities.InOutURL) (entities.InOutURL, error) {
			return entities.InOutURL{}, apperor.ErrRepoNotFound
		},
	}

	usecase := postgres.NewUseCase(mockRepo)
	ctx := context.Background()

	getURL := entities.InOutURL{URL: "nonexistent"}
	_, err := usecase.GetShortURL(ctx, &getURL)

	if err == nil {
		t.Error("Expected error for non-existent URL, got nil")
	}
}

func TestPostgresUsecase_CreateDuplicate(t *testing.T) {

	mockRepo := &MockRepository{
		CreateFunc: func(ctx context.Context, url entities.URLsStruct) error {
			_ = url.ShortURL
			return nil
		},
		GetFunc: func(ctx context.Context, url entities.InOutURL) (entities.InOutURL, error) {
			return entities.InOutURL{URL: "https://ozon.ru"}, nil
		},
	}

	usecase := postgres.NewUseCase(mockRepo)
	ctx := context.Background()

	postURL := entities.InOutURL{URL: "https://ozon.ru"}

	shortURL_1, err := usecase.CreateShortURL(ctx, &postURL)
	if err != nil {
		t.Fatalf("First CreateShortURL failed: %v", err)
	}

	shortURL_2, err := usecase.CreateShortURL(ctx, &postURL)
	if err != nil {
		t.Fatalf("Second CreateShortURL failed: %v", err)
	}

	if shortURL_1.URL != shortURL_2.URL {
		t.Errorf("Expected same short URL, got %s and %s", shortURL_1.URL, shortURL_2.URL)
	}
}
