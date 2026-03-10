package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"url-shortener-ozon/pkg/utils"

	"github.com/stretchr/testify/assert"

	apperor "url-shortener-ozon/internal/apperror"
	"url-shortener-ozon/internal/domain/entities"
	"url-shortener-ozon/internal/domain/usecase"
	memrepo "url-shortener-ozon/internal/repository/url/memory"
)

type MockRepository struct {
	CreateFunc func(ctx context.Context, url entities.URLsStruct) error
	GetFunc    func(ctx context.Context, url entities.RequestData) (entities.ResponseData, error)
}

func (m *MockRepository) CreateShortPath(ctx context.Context, url entities.URLsStruct) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, url)
	}
	return nil
}

func (m *MockRepository) GetOriginalURLByShortPath(ctx context.Context, url entities.RequestData) (entities.ResponseData, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, url)
	}
	return entities.ResponseData{}, nil
}

// Post запросы

// Создаём ссылку, получаем ответ, нет ошибки
func TestCreateShortPath_WithPostgres_Success(t *testing.T) {

	expectedOriginalURL := "https://ozon.ru"

	// Вычисляем ожидаемый shortPath для salt = 0
	expectedShortPath := utils.GenerateShortPath(expectedOriginalURL, 0)

	// Флаг для отслеживания вызова CreateShortPath
	createCalled := false

	mockRepo := &MockRepository{
		// Мокаем GetOriginalURLByShortPath - возвращаем "не найдено"
		GetFunc: func(ctx context.Context, url entities.RequestData) (entities.ResponseData, error) {
			// Проверяем, что запрашивается правильный shortPath
			assert.Equal(t, expectedShortPath, url.URL)
			return entities.ResponseData{}, apperor.ErrRepoNotFound
		},
		// Мокаем CreateShortPath - должен быть вызван один раз
		CreateFunc: func(ctx context.Context, url entities.URLsStruct) error {
			createCalled = true
			assert.Equal(t, expectedOriginalURL, url.OriginalURL)
			assert.Equal(t, expectedShortPath, url.ShortPath)
			assert.NotEmpty(t, url.ShortPath)
			return nil
		},
	}

	uc := usecase.NewUseCase(mockRepo)

	requestData := &entities.RequestData{URL: expectedOriginalURL}
	ctx := context.Background()
	response, err := uc.CreateShortPath(ctx, requestData)

	assert.NoError(t, err)
	assert.Equal(t, expectedShortPath, response.URL)
	assert.True(t, createCalled, "CreateShortPath should be called")
}

func TestCreateShortPath_WithMemory_Success(t *testing.T) {
	repo := memrepo.NewMemoryRepository()
	uc := usecase.NewUseCase(repo)
	ctx := context.Background()

	testURL := "https://ozon.com"

	// Создаем короткую ссылку
	input := entities.RequestData{URL: testURL}
	shortPath, err := uc.CreateShortPath(ctx, &input)

	assert.NoError(t, err, "CreateShortPath should not return error")
	assert.NotEmpty(t, shortPath.URL, "Short URL should not be empty")
}

// Создаём ссылки для разных URL. Должны получиться разные короткие пути
func TestCreateShortURL_WithPostgres_DifferentShortPathsForDifferentURLs(t *testing.T) {

	shortToOriginalMap := make(map[string]string) // map[shortPath]originalURL
	getCalls := make(map[string]bool)             // отслеживаем вызовы Get

	mockRepo := &MockRepository{
		// Мокаем GetFunc, чтобы возвращать "не найдено" для всех запросов
		GetFunc: func(ctx context.Context, url entities.RequestData) (entities.ResponseData, error) {
			getCalls[url.URL] = true

			// Всегда возвращаем "не найдено", чтобы usecase пытался создать новую ссылку
			// Но проверяем, что запрашивается именно тот shortPath, который мы ожидаем
			assert.NotEmpty(t, url.URL, "Requested short path should not be empty")

			// Возвращаем ErrRepoNotFound, чтобы usecase понял, что ссылка уникальна
			return entities.ResponseData{}, apperor.ErrRepoNotFound
		},

		CreateFunc: func(ctx context.Context, url entities.URLsStruct) error {
			// Проверяем, что для этого короткого пути еще не сохраняли оригинальный URL
			if existingOriginal, exists := shortToOriginalMap[url.ShortPath]; exists {
				// Если коллизия, проверяем что это тот же оригинальный URL
				assert.Equal(t, existingOriginal, url.OriginalURL, "Same short path should map to same original URL")
			} else {
				// Сохраняем соответствие короткого пути оригинальному URL
				shortToOriginalMap[url.ShortPath] = url.OriginalURL
			}
			assert.NotEmpty(t, url.ShortPath, "ShortPath should not be empty")
			return nil
		},
	}

	uc := usecase.NewUseCase(mockRepo)

	// Создаем ссылки для разных URL
	testURLs := []string{
		"https://ozon.ru",
		"https://ya.ru",
		"https://github.com",
	}

	var shortPaths []string

	for _, testURL := range testURLs {
		originalURL := entities.RequestData{URL: testURL}
		ctx := context.Background()
		shortPath, err := uc.CreateShortPath(ctx, &originalURL)

		assert.NoError(t, err, "CreateShortPath should not return error for URL: %s", testURL)
		assert.NotEmpty(t, shortPath.URL, "Short URL should not be empty for URL: %s", testURL)

		shortPaths = append(shortPaths, shortPath.URL)
	}

	// Проверяем, что GetFunc вызывался для каждого созданного shortPath
	for _, shortPath := range shortPaths {
		assert.True(t, getCalls[shortPath], "GetFunc should have been called for short path: %s", shortPath)
	}

	// Проверяем, что все короткие пути разные
	for i := 0; i < len(shortPaths); i++ {
		for j := i + 1; j < len(shortPaths); j++ {
			assert.NotEqual(t, shortPaths[i], shortPaths[j],
				"Short paths for different URLs should be unique: %s and %s are equal",
				shortPaths[i], shortPaths[j])
		}
	}

	// Проверяем, что для каждого короткого пути сохранен свой оригинальный URL
	assert.Equal(t, len(testURLs), len(shortToOriginalMap), "Should have created unique short paths for each URL")

	// Проверяем, что все оригинальные URL присутствуют в map
	for _, testURL := range testURLs {
		found := false
		for _, originalURL := range shortToOriginalMap {
			if originalURL == testURL {
				found = true
				break
			}
		}
		assert.True(t, found, "Original URL %s should be mapped from some short path", testURL)
	}
}

func TestCreateShortURL_WithMemory_DifferentShortPathsForDifferentURLs(t *testing.T) {
	repo := memrepo.NewMemoryRepository()
	uc := usecase.NewUseCase(repo)
	ctx := context.Background()

	// Создаем ссылки для разных URL
	testURLs := []string{
		"https://ozon.ru",
		"https://ya.ru",
		"https://github.com",
	}

	var shortPaths []string
	shortToOriginalMap := make(map[string]string) // map[shortPath]originalURL

	for _, testURL := range testURLs {
		originalURL := entities.RequestData{URL: testURL}
		shortPath, err := uc.CreateShortPath(ctx, &originalURL)

		assert.NoError(t, err, "CreateShortPath should not return error for URL: %s", testURL)
		assert.NotEmpty(t, shortPath.URL, "Short URL should not be empty for URL: %s", testURL)

		shortPaths = append(shortPaths, shortPath.URL)
		shortToOriginalMap[shortPath.URL] = testURL
	}

	// Проверяем, что все короткие пути разные
	for i := 0; i < len(shortPaths); i++ {
		for j := i + 1; j < len(shortPaths); j++ {
			assert.NotEqual(t, shortPaths[i], shortPaths[j],
				"Short paths for different URLs should be unique. Found duplicate: %s",
				shortPaths[i])
		}
	}

	// Проверяем, что для каждого короткого пути сохранен свой оригинальный URL
	for shortPath, originalURL := range shortToOriginalMap {
		// Получаем оригинальный URL по короткому пути
		inputPath := entities.RequestData{URL: shortPath}
		inputURL, err := uc.GetOriginalURLByShortPath(ctx, &inputPath)

		assert.NoError(t, err, "Should be able to retrieve URL by short path: %s", shortPath)
		assert.Equal(t, originalURL, inputURL.URL,
			"Short path %s should map back to original URL %s", shortPath, originalURL)
	}

	// Проверяем, что для одного и того же URL всегда получаем один и тот же короткий путь
	for shortPath, originalURL := range shortToOriginalMap {
		inputURL := entities.RequestData{URL: originalURL}
		inputPath, err := uc.CreateShortPath(ctx, &inputURL)

		assert.NoError(t, err, "CreateShortPath for existing URL should not return error")
		assert.Equal(t, shortPath, inputPath.URL,
			"Same original URL %s should always return same short path %s",
			originalURL, shortPath)
	}
}

// Ошибка репозитория при создании
func TestCreateShortURL_WithPostgres_CreateRepositoryError(t *testing.T) {

	expectedErr := errors.New("database connection error") // Используем реальную ошибку, не ErrRepoNotFound

	mockRepo := &MockRepository{
		// Мокаем GetFunc, чтобы он возвращал "не найдено"
		GetFunc: func(ctx context.Context, url entities.RequestData) (entities.ResponseData, error) {
			// Возвращаем "не найдено", чтобы usecase попытался создать
			return entities.ResponseData{}, apperor.ErrRepoNotFound
		},
		// CreateFunc возвращает ошибку
		CreateFunc: func(ctx context.Context, url entities.URLsStruct) error {
			return expectedErr
		},
	}

	uc := usecase.NewUseCase(mockRepo)

	originalURL := entities.RequestData{URL: "https://ozon.ru"}
	ctx := context.Background()
	_, err := uc.CreateShortPath(ctx, &originalURL)

	assert.Error(t, err, "Expected error from repository")
	assert.Equal(t, expectedErr, err, "Error should match expected error")
}

// Get запросы

// Передаём путь, получаем ссылку, нет ошибки
func TestGetOriginalURL_WithPostgres_Success(t *testing.T) {
	mockRepo := &MockRepository{
		GetFunc: func(ctx context.Context, url entities.RequestData) (entities.ResponseData, error) {
			assert.Equal(t, "short-url", url.URL, "Request URL should match")
			return entities.ResponseData{URL: "https://ozon.ru"}, nil
		},
	}

	uc := usecase.NewUseCase(mockRepo)
	ctx := context.Background()

	// Получаем оригинальную ссылку
	shortPath := entities.RequestData{URL: "short-url"}
	originalURL, err := uc.GetOriginalURLByShortPath(ctx, &shortPath)

	assert.NoError(t, err, "GetOriginalURLByShortPath should not return error")
	assert.Equal(t, "https://ozon.ru", originalURL.URL, "Original URL should match expected")
}

func TestGetOriginalURL_WithMemory_Success(t *testing.T) {
	repo := memrepo.NewMemoryRepository()
	uc := usecase.NewUseCase(repo)
	ctx := context.Background()

	// Сначала создаем запись
	testURL := "https://ozon.ru"
	createPath := entities.RequestData{URL: testURL}
	shortPath, err := uc.CreateShortPath(ctx, &createPath)
	assert.NoError(t, err, "CreateShortPath should not return error")
	assert.NotEmpty(t, shortPath.URL, "Short URL should not be empty")

	// Получаем оригинальную ссылку
	inputPath := entities.RequestData{URL: shortPath.URL}
	originalURL, err := uc.GetOriginalURLByShortPath(ctx, &inputPath)

	assert.NoError(t, err, "GetOriginalURLByShortPath should not return error")
	assert.Equal(t, testURL, originalURL.URL, "Original URL should match expected")
}

// Попытка получить ответ без возврата
func TestGetOriginalURL_WithPostgres_NotFound(t *testing.T) {
	mockRepo := &MockRepository{
		GetFunc: func(ctx context.Context, url entities.RequestData) (entities.ResponseData, error) {
			assert.Equal(t, "nonexistent", url.URL, "Request URL should match")
			return entities.ResponseData{}, apperor.ErrRepoNotFound
		},
	}

	uc := usecase.NewUseCase(mockRepo)
	ctx := context.Background()

	// Пытаемся получить несуществующую ссылку
	shortPath := entities.RequestData{URL: "nonexistent"}
	_, err := uc.GetOriginalURLByShortPath(ctx, &shortPath)

	assert.Error(t, err, "Expected error for non-existent URL")
	assert.Equal(t, apperor.ErrRepoNotFound, err, "Error should be ErrRepoNotFound")
}

func TestGetOriginalURL_WithMemory_NotFound(t *testing.T) {
	repo := memrepo.NewMemoryRepository()
	uc := usecase.NewUseCase(repo)
	ctx := context.Background()

	// Пытаемся получить несуществующую ссылку
	shortPath := entities.RequestData{URL: "nonexistent"}
	_, err := uc.GetOriginalURLByShortPath(ctx, &shortPath)

	assert.Error(t, err, "Expected error for non-existent URL")
	assert.Equal(t, apperor.ErrRepoNotFound, err, "Error should be ErrRepoNotFound")
}

// Тест на не включённый репозиторий
func TestGetOriginalURL_WithPostgres_QueryExecutionError(t *testing.T) {
	expectedErr := fmt.Errorf("executing query error: closed pool")

	mockRepo := &MockRepository{
		GetFunc: func(ctx context.Context, url entities.RequestData) (entities.ResponseData, error) {
			assert.Equal(t, "some-url", url.URL, "Request URL should match")
			return entities.ResponseData{}, expectedErr
		},
	}

	uc := usecase.NewUseCase(mockRepo)
	ctx := context.Background()

	shortPath := entities.RequestData{URL: "some-url"}
	_, err := uc.GetOriginalURLByShortPath(ctx, &shortPath)

	assert.Error(t, err, "Expected error from repository")
	assert.Equal(t, expectedErr.Error(), err.Error(), "Error message should match expected")
	assert.Contains(t, err.Error(), "closed pool", "Error should indicate connection pool is closed")
}

// Тест на ошибки при парсинге результата
func TestGetOriginalURL_WithPostgres_CollectRowError(t *testing.T) {
	expectedErr := fmt.Errorf("selecting url from database error: some db error")

	mockRepo := &MockRepository{
		GetFunc: func(ctx context.Context, url entities.RequestData) (entities.ResponseData, error) {
			assert.Equal(t, "some-url", url.URL, "Request URL should match")
			return entities.ResponseData{}, expectedErr
		},
	}

	uc := usecase.NewUseCase(mockRepo)
	ctx := context.Background()

	shortPath := entities.RequestData{URL: "some-url"}
	_, err := uc.GetOriginalURLByShortPath(ctx, &shortPath)

	assert.Error(t, err, "Expected error from repository")
	assert.Equal(t, expectedErr.Error(), err.Error(), "Error message should match expected")
	assert.Contains(t, err.Error(), "selecting url from database", "Error should indicate issue with collecting row")
}

// Тест на проверку обработки контекстных ошибок
func TestGetOriginalURL_WithPostgres_ContextCancelled(t *testing.T) {
	expectedErr := context.Canceled

	mockRepo := &MockRepository{
		GetFunc: func(ctx context.Context, url entities.RequestData) (entities.ResponseData, error) {
			assert.Equal(t, "some-url", url.URL, "Request URL should match")
			return entities.ResponseData{}, expectedErr
		},
	}

	uc := usecase.NewUseCase(mockRepo)
	ctx := context.Background()

	shortPath := entities.RequestData{URL: "some-url"}
	_, err := uc.GetOriginalURLByShortPath(ctx, &shortPath)

	assert.Error(t, err, "Expected context cancellation error")
	assert.Equal(t, expectedErr, err, "Error should be context.Canceled")
}
