package memory_test

import (
	"context"
	"testing"

	"url-shortener-ozon/internal/domain/entities"
	"url-shortener-ozon/internal/domain/usecase/memory"
	memrepo "url-shortener-ozon/internal/repository/url/memory"
)

func TestMemoryUsecase_CreateAndGet(t *testing.T) {
	// Создаем репозиторий и usecase
	repo := memrepo.NewMemoryRepository()
	uc := memory.NewUseCase(repo)

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
	uc := memory.NewUseCase(repo)

	ctx := context.Background()

	// Пытаемся получить несуществующую ссылку
	input := entities.InOutURL{URL: "nonexistent"}
	_, err := uc.GetShortURL(ctx, &input)

	if err == nil {
		t.Fatal("Expected error for non-existent URL, got nil")
	}
}
