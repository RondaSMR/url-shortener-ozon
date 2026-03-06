package postgres_test

import (
	"context"
	"os"
	"testing"

	"url-shortener-ozon/internal/domain/entities"
	"url-shortener-ozon/internal/domain/usecase/postgres"
	pgrepo "url-shortener-ozon/internal/repository/url/postgres"
	"url-shortener-ozon/pkg/connectors/pgconnector"
)

func TestPostgresUsecase_CreateAndGet(t *testing.T) {
	// Пропускаем тест, если нет переменных окружения
	host := os.Getenv("STORAGE_HOST")
	if host == "" {
		t.Skip("Skipping postgres test: STORAGE_HOST not set")
	}

	// Подключаемся к БД
	cfg, err := pgconnector.CreateConfig(&pgconnector.ConnectionConfig{
		Host:     host,
		Port:     os.Getenv("STORAGE_PORT"),
		User:     os.Getenv("STORAGE_PG_USER"),
		Password: os.Getenv("STORAGE_PASS"),
		DbName:   os.Getenv("STORAGE_DB"),
		SslMode:  "disable",
	}, nil)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	conn, err := pgconnector.NewPgConnector(cfg, 5, 5)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer conn.CloseConnection()

	// Создаем репозиторий и usecase
	repo := pgrepo.NewPostgresRepository(conn)
	uc := postgres.NewUseCase(repo)

	ctx := context.Background()

	// Тестовые данные
	testURL := "https://ozon.ru"

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
}

func TestPostgresUsecase_CreateDuplicate(t *testing.T) {
	// Пропускаем тест, если нет переменных окружения
	host := os.Getenv("STORAGE_HOST")
	if host == "" {
		t.Skip("Skipping postgres test: STORAGE_HOST not set")
	}

	// Подключаемся к БД
	cfg, err := pgconnector.CreateConfig(&pgconnector.ConnectionConfig{
		Host:     host,
		Port:     os.Getenv("STORAGE_PORT"),
		User:     os.Getenv("STORAGE_PG_USER"),
		Password: os.Getenv("STORAGE_PASS"),
		DbName:   os.Getenv("STORAGE_DB"),
		SslMode:  "disable",
	}, nil)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	conn, err := pgconnector.NewPgConnector(cfg, 5, 5)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer conn.CloseConnection()

	// Создаем репозиторий и usecase
	repo := pgrepo.NewPostgresRepository(conn)
	uc := postgres.NewUseCase(repo)

	ctx := context.Background()

	// Создаем ссылку первый раз
	testURL := "https://ozon.ru"
	input := entities.InOutURL{URL: testURL}

	short1, err := uc.CreateShortURL(ctx, &input)
	if err != nil {
		t.Fatalf("First CreateShortURL failed: %v", err)
	}

	// Создаем ту же ссылку второй раз (должен вернуть ту же короткую)
	short2, err := uc.CreateShortURL(ctx, &input)
	if err != nil {
		t.Fatalf("Second CreateShortURL failed: %v", err)
	}

	if short1.URL != short2.URL {
		t.Fatalf("Expected same short URL for duplicate, got %q and %q", short1.URL, short2.URL)
	}
}
