package config

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

type HttpServer struct {
	Address string `env:"HTTP_SERVER_ADDRESS" yaml:"address"`
	User    string `env:"HTTP_SERVER_USER" yaml:"user"`
	Pass    string `env:"HTTP_SERVER_PASS" yaml:"pass"`
}

type AppConfig struct {
	PathConfig  string     `env:"PATH_CONFIG"`
	ServiceName string     `env:"SERVICE_NAME" yaml:"serviceName"`
	Debug       bool       `env:"DEBUG" yaml:"debug"`
	StorageMode string     `env:"STORAGE_MODE" envDefault:"postgres" yaml:"storageMode"`
	PGStorage   pgStorage  `yaml:"pgStorage"`
	HTTPServer  HttpServer `yaml:"http_server"`
}

type pgStorage struct {
	Host string `env:"STORAGE_HOST" yaml:"host"`
	Port int    `env:"STORAGE_PORT" yaml:"port"`
	User string `env:"STORAGE_PG_USER" yaml:"user"`
	Pass string `env:"STORAGE_PASS" yaml:"pass"`
	DB   string `env:"STORAGE_DB"   yaml:"db"`
}

// ReadYamlConfig считывает файл конфигурации YAML и сохраняет его содержимое в структуре AppConfig
func (c *AppConfig) ReadYamlConfig(pathFile string) error {
	open, err := os.Open(pathFile)
	if err != nil {
		return err
	}
	defer func(open *os.File) {
		err = open.Close()
		if err != nil {
			log.Println(err)
		}
	}(open)
	if err = yaml.NewDecoder(open).Decode(c); err != nil { //nolint:typecheck
		return err
	}
	return nil
}

// ReadEnvConfig считывает переменные окружения и сохраняет их содержимое в структуре AppConfig
func (c *AppConfig) ReadEnvConfig() error {
	if err := env.Parse(c); err != nil { //nolint:typecheck
		return err
	}
	return nil
}

// Validate проверяет структуру AppConfig
func (c *AppConfig) Validate() error {
	if c.HTTPServer.Address == "" {
		return fmt.Errorf("address service is not set")
	}

	mode := strings.ToLower(strings.TrimSpace(c.StorageMode))
	switch mode {
	case "", "memory":
		c.StorageMode = "memory"
		return nil
	case "db", "postgres", "pg":
		if c.PGStorage.Host == "" {
			return fmt.Errorf("pg repository host is not set")
		}
		if c.PGStorage.Port <= 0 {
			return fmt.Errorf("pg repository port is not set")
		}
		if c.PGStorage.User == "" {
			return fmt.Errorf("pg repository user is not set")
		}
		if c.PGStorage.Pass == "" {
			return fmt.Errorf("pg repository password is not set")
		}
		if c.PGStorage.DB == "" {
			return fmt.Errorf("pg repository db name is not set")
		}
		c.StorageMode = "db"
		return nil
	default:
		return fmt.Errorf("unknown STORAGE_MODE=%q (expected: memory|db)", c.StorageMode)
	}
}
