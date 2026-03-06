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
	"url-shortener-ozon/internal/controller/http/v1/urlshortener"
	mem2 "url-shortener-ozon/internal/domain/usecase/memory"
	pg2 "url-shortener-ozon/internal/domain/usecase/postgres"
	mem3 "url-shortener-ozon/internal/repository/url/memory"
	pg3 "url-shortener-ozon/internal/repository/url/postgres"
	"url-shortener-ozon/pkg/config"
	"url-shortener-ozon/pkg/connectors/pgconnector"
)

const (
	timeOutShutdownService = time.Duration(5) * time.Second
	timeReadHeader         = time.Duration(5) * time.Second
)

func NewAppPostgres(config config.AppConfig) error {

	// Инициализация PostgreSQL
	connectorConfig, err := pgconnector.CreateConfig(&pgconnector.ConnectionConfig{
		Host:     config.PGStorage.Host,
		Port:     fmt.Sprint(config.PGStorage.Port),
		User:     config.PGStorage.User,
		Password: config.PGStorage.Pass,
		DbName:   config.PGStorage.DB,
		SslMode:  "disable",
	},
		nil)

	pgStorage, err := pgconnector.NewPgConnector(
		connectorConfig,
		10*time.Second,
		10*time.Second,
	)
	if err != nil {
		return fmt.Errorf("initialize pg storage: %w", err)
	}
	defer func() {
		pgStorage.CloseConnection()
	}()

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
		pg2.NewUseCase(pg3.NewPostgresRepository(pgStorage)),
		config.HTTPServer,
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
	return nil
}

func NewAppMemory(config config.AppConfig) error {

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
		mem2.NewUseCase(mem3.NewMemoryRepository()),
		config.HTTPServer,
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
	return nil
}

// Функция инициализации роутера
func routersInit(
	router *gin.Engine,
	usecase urlshortener.Usecase,
	srv config.HttpServer,
) {
	urlshortener.Router(router.Group("url-shortener-ozon"), usecase)
}
