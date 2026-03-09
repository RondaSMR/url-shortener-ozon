package usecase_test

import (
	"context"
	"testing"
	apperor "url-shortener-ozon/internal/apperror"
	"url-shortener-ozon/internal/domain/usecase"
	memrepo "url-shortener-ozon/internal/repository/url/memory"

	"url-shortener-ozon/internal/domain/entities"
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

	uc := usecase.NewUseCase(mockRepo)
	ctx := context.Background()

	postURL := entities.InOutURL{URL: "https://ozon.ru"}
	shortURL, err := uc.CreateShortURL(ctx, &postURL)
	if err != nil {
		t.Fatalf("CreateShortURL failed: %v", err)
	}
	if shortURL.URL == "" {
		t.Fatal("CreateShortURL returned empty URL")
	}

	getURL := entities.InOutURL{URL: shortURL.URL}
	originalURL, err := uc.GetShortURL(ctx, &getURL)
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

	uc := usecase.NewUseCase(mockRepo)
	ctx := context.Background()

	getURL := entities.InOutURL{URL: "nonexistent"}
	_, err := uc.GetShortURL(ctx, &getURL)

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

	uc := usecase.NewUseCase(mockRepo)
	ctx := context.Background()

	postURL := entities.InOutURL{URL: "https://ozon.ru"}

	shortURL_1, err := uc.CreateShortURL(ctx, &postURL)
	if err != nil {
		t.Fatalf("First CreateShortURL failed: %v", err)
	}

	shortURL_2, err := uc.CreateShortURL(ctx, &postURL)
	if err != nil {
		t.Fatalf("Second CreateShortURL failed: %v", err)
	}

	if shortURL_1.URL != shortURL_2.URL {
		t.Errorf("Expected same short URL, got %s and %s", shortURL_1.URL, shortURL_2.URL)
	}
}

// ------

func TestMemoryUsecase_CreateAndGet(t *testing.T) {
	// Создаем репозиторий и usecase
	repo := memrepo.NewMemoryRepository()
	uc := usecase.NewUseCase(repo)

	ctx := context.Background()

	// Тестовые данные
	testURLs := []string{
		"https://ozon.com",
		"https://google.com",
		"https://github.com",
	}

	for _, testURL := range testURLs {
		t.Run(testURL, func(t *testing.T) {
			// Создаем короткую ссылку
			input := entities.InOutURL{URL: testURL}
			short, err := uc.CreateShortURL(ctx, &input)

			if err != nil {
				t.Fatalf("CreateShortURL failed: %v", err)
			}
			if short.URL == "" {
				t.Fatal("CreateShortURL returned empty URL")
			}

			// Получаем оригинальную ссылку
			getInput := entities.InOutURL{URL: short.URL}
			original, err := uc.GetShortURL(ctx, &getInput)

			if err != nil {
				t.Fatalf("GetShortURL failed: %v", err)
			}
			if original.URL != testURL {
				t.Fatalf("Expected %q, got %q", testURL, original.URL)
			}
		})
	}
}

func TestMemoryUsecase_GetNonExistent(t *testing.T) {
	repo := memrepo.NewMemoryRepository()
	uc := usecase.NewUseCase(repo)

	ctx := context.Background()

	// Пытаемся получить несуществующую ссылку
	input := entities.InOutURL{URL: "nonexistent"}
	_, err := uc.GetShortURL(ctx, &input)

	if err == nil {
		t.Fatal("Expected error for non-existent URL, got nil")
	}
}
