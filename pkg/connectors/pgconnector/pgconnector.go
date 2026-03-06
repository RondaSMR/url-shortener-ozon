package pgconnector

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgConnector interface {
	GetPool() *pgxpool.Pool
	CloseConnection()
	Ping() error
}

type Connector struct {
	poolMu      sync.Mutex
	pool        *pgxpool.Pool
	pingTimeout time.Duration
}

type ConnectionConfig struct {
	Host     string `form:"host"`
	Port     string `form:"port"`
	DbName   string `form:"dbname"`
	User     string `form:"user"`
	Password string `form:"password"`
	SslMode  string `form:"sslmode"`
}

type PoolConfig struct {
	MaxOpenConns          int           `form:"pool_max_open_conns"`
	MaxIdleConns          int           `form:"pool_max_idle_conns"`
	MaxConnLifetime       time.Duration `form:"pool_max_conn_lifetime"`
	MaxConnIdleTime       time.Duration `form:"pool_max_conn_idle_time"`
	HealthCheckPeriod     time.Duration `form:"pool_health_check_period"`
	MaxConnLifetimeJitter time.Duration `form:"pool_max_conn_lifetime_jitter"`
}

func NewPgConnector(
	config *pgxpool.Config,
	pingTimeout time.Duration,
	connectTimeout time.Duration,
) (*Connector, error) {
	if config == nil {
		return nil, fmt.Errorf("config is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	connector := Connector{
		pool:        pool,
		pingTimeout: pingTimeout,
	}
	return &connector, connector.Ping()
}

// CreateConfig создает pgxpool.Config из ConnectionConfig и PoolConfig (если существует)
//
// Вернёт ошибку если connectionConfig == nil
func CreateConfig(
	connectionConfig *ConnectionConfig,
	poolConfig *PoolConfig,
) (*pgxpool.Config, error) {
	if connectionConfig == nil {
		return nil, fmt.Errorf("connection config is nil")
	}

	var connString bytes.Buffer

	// Парсит конфиг подключения
	reflectValue := reflect.ValueOf(*connectionConfig)
	for i := range reflectValue.Type().NumField() {
		field := reflectValue.Type().Field(i)
		value := reflectValue.Field(i)
		tag := field.Tag.Get("form")

		if value.String() == "" {
			return nil, fmt.Errorf("empty field: %s", tag)
		}

		connString.WriteString(
			fmt.Sprintf("%s=%s ",
				tag,
				value.String(),
			),
		)
	}

	// Парсит poolConfig если существует
	if poolConfig != nil {
		reflectValue = reflect.ValueOf(poolConfig)
		for i := range reflectValue.Type().NumField() {
			field := reflectValue.Type().Field(i)
			value := reflectValue.Field(i)
			tag := field.Tag.Get("form")

			if value.String() == "" {
				continue
			}

			connString.WriteString(
				fmt.Sprintf("%s=%s ",
					tag,
					value.String(),
				),
			)
		}
	}

	return pgxpool.ParseConfig(connString.String())
}

func (p *Connector) GetPool() *pgxpool.Pool {
	p.poolMu.Lock()
	defer p.poolMu.Unlock()
	return p.pool
}

func (p *Connector) CloseConnection() {
	p.pool.Close()
}

func (p *Connector) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), p.pingTimeout)
	defer cancel()

	return p.pool.Ping(ctx)
}
