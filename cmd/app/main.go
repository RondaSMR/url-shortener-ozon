package main

import (
	"go.uber.org/zap"
	"log"
	"strings"
	"url-shortener-ozon/internal/app"
	"url-shortener-ozon/pkg/config"
	"url-shortener-ozon/pkg/utils"
)

func main() {
	cfg := initAppConfig()
	initStorageMode(cfg)
}

// Инициализация конфигурации
func initAppConfig() config.AppConfig {
	var cfg config.AppConfig
	// Инициализация логгера
	logger, err := utils.CreateLogger(zap.InfoLevel)
	if err != nil {
		log.Fatalf("creating logger failed: %v", err)
	}
	if err = cfg.ReadEnvConfig(); err != nil {
		logger.Fatal("reading environment variables failed", zap.Error(err))
	}
	if cfg.PathConfig != "" {
		if err = cfg.ReadYamlConfig(cfg.PathConfig); err != nil {
			logger.Fatal("reading config failed", zap.Error(err))
		}
	}

	if err = cfg.Validate(); err != nil {
		logger.Fatal("validating config failed", zap.Error(err))
	}
	if cfg.Debug {
		logger.Warn("application is running in debug mode")
		logger, err = utils.CreateLogger(zap.DebugLevel)
		if err != nil {
			log.Fatalf("failed to create logger: %s", err)
		}
	}
	zap.ReplaceGlobals(logger)
	return cfg
}

func initStorageMode(cfg config.AppConfig) {
	// Выбор режима записи укороченных ссылок
	switch strings.ToLower(strings.TrimSpace(cfg.StorageMode)) {
	case "memory":
		if err := app.NewAppMemory(cfg); err != nil {
			zap.L().Fatal("application failed", zap.Error(err))
			return
		}
	case "db":
		if err := app.NewAppPostgres(cfg); err != nil {
			zap.L().Fatal("application failed", zap.Error(err))
			return
		}
	default:
		zap.L().Fatal("unknown STORAGE_MODE", zap.String("storage_mode", cfg.StorageMode))
	}
}
