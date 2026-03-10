package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener-ozon/internal/controller/http/v1/url_shortener"
	"url-shortener-ozon/internal/domain/usecase"
	"url-shortener-ozon/internal/repository/url/memory"
	"url-shortener-ozon/internal/repository/url/postgres"
	"url-shortener-ozon/pkg/config"
)

const (
	timeOutShutdownService = time.Duration(5) * time.Second
	timeReadHeader         = time.Duration(5) * time.Second
)

// NewAppPostgres инициализирует работу в режиме сохранения в базу данных Postgres
func NewAppPostgres(config config.AppConfig) error {

	repo, err := postgres.NewPostgresRepository(config.PGStorage)
	if err != nil {
		return fmt.Errorf("initialize pg storage: %w", err)
	}
	defer repo.Close()

	routersInit2(config, repo)
	return nil
}

// NewAppMemory инициализирует работу в режиме in-memory сохранения
func NewAppMemory(config config.AppConfig) error {
	routersInit2(config, memory.NewMemoryRepository())
	return nil
}

// Функция включения роутера
func routersInit2(config config.AppConfig, repository usecase.Repository) {

	// Настройка роутера
	router := gin.New()

	if config.Debug {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	router.Use(otelgin.Middleware("url-service"))
	routersInit(
		router,
		usecase.NewUseCase(repository),
	)

	srv := &http.Server{
		Addr:        config.HTTPServer.Address,
		Handler:     router,
		ReadTimeout: timeReadHeader,
	}

	go func() {
		zap.L().Info("Server is starting", zap.String("address", "http://"+config.HTTPServer.Address))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), timeOutShutdownService)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("emergency shutdown http server", zap.Error(err))
	}
	zap.L().Info("http server shutdown")
}

// Функция инициализации роутера HTTP-запросов
func routersInit(
	router *gin.Engine,
	usecase url_shortener.Usecase,
) {
	url_shortener.Router(router.Group("url-shortener-ozon"), usecase)
}
