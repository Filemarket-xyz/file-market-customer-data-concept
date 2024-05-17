package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

var Mode string

type (
	Config struct {
		Server       *ServerConfig
		Handler      *HandlerConfig
		Service      *ServiceConfig
		Redis        *RedisConfig
		Postgres     *PostgresConfig
		TokenManager *TokenManagerConfig
		Locale       int64
		Dune         *DuneConfig
	}

	RedisConfig struct {
		Host          string
		Port          string
		Password      string
		UserTokensTTL time.Duration
	}

	PostgresConfig struct {
		Host     string
		User     string
		Password string
		DBName   string
		Port     int
	}

	ServerConfig struct {
		Port           int
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
		MaxHeaderBytes int
	}

	HandlerConfig struct {
		RequestTimeout  time.Duration
		RegisterTimeout time.Duration
		PayTimeout      time.Duration
		TgKey           string
	}

	ServiceConfig struct {
		AccessTokenTTL    time.Duration
		RefreshTokenTTL   time.Duration
		SignatureDuration time.Duration
		RpcUrl            string
		Mode              string
		Wallet            common.Address
		DatasetCost       decimal.Decimal
		Periods           *Periods
	}

	Periods struct {
		ListenerPeriod time.Duration
		JwtCheckPeriod time.Duration
	}

	TokenManagerConfig struct {
		SigningKey string
	}

	DuneConfig struct {
		ApiKey string
	}
)

func Init(configPath string) (*Config, error) {
	jsonCfg := viper.New()
	jsonCfg.AddConfigPath(filepath.Dir(configPath))
	jsonCfg.SetConfigName(filepath.Base(configPath))

	if err := jsonCfg.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config/Init/jsonCfg.ReadInConfig: %w", err)
	}

	envCfg := viper.New()
	envCfg.SetConfigFile(".env")

	if err := envCfg.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config/Init/envCfg.ReadInConfig: %w", err)
	}

	walletAddress := jsonCfg.GetString("service.walletAddress")
	wallet := common.HexToAddress(walletAddress)

	datasetCostStr := jsonCfg.GetString("service.datasetCost")
	datasetCost, err := decimal.NewFromString(datasetCostStr)
	if err != nil {
		return nil, fmt.Errorf("config/Init/decimal.NewFromString: %w", err)
	}

	return &Config{
		Locale: jsonCfg.GetInt64("locale"),
		Server: &ServerConfig{
			Port:           jsonCfg.GetInt("server.port"),
			ReadTimeout:    jsonCfg.GetDuration("server.readTimeout"),
			WriteTimeout:   jsonCfg.GetDuration("server.writeTimeout"),
			MaxHeaderBytes: jsonCfg.GetInt("server.maxHeaderBytes"),
		},
		Handler: &HandlerConfig{
			RequestTimeout:  jsonCfg.GetDuration("handler.requestTimeout"),
			RegisterTimeout: jsonCfg.GetDuration("handler.registerTimeout"),
			PayTimeout:      jsonCfg.GetDuration("handler.payTimeout"),
			TgKey:           envCfg.GetString("TG_KEY"),
		},
		Service: &ServiceConfig{
			AccessTokenTTL:    jsonCfg.GetDuration("service.accessTTL"),
			RefreshTokenTTL:   jsonCfg.GetDuration("service.refreshTTL"),
			SignatureDuration: jsonCfg.GetDuration("service.signature_duration"),
			Mode:              jsonCfg.GetString("service.mode"),
			RpcUrl:            envCfg.GetString("RPC_URL"),
			Periods: &Periods{
				ListenerPeriod: jsonCfg.GetDuration("service.periods.listener"),
				JwtCheckPeriod: jsonCfg.GetDuration("service.periods.jwtCheckPeriod"),
			},
			Wallet:      wallet,
			DatasetCost: datasetCost,
		},
		Redis: &RedisConfig{
			Host:          envCfg.GetString("REDIS_HOST"),
			Port:          envCfg.GetString("REDIS_PORT"),
			Password:      envCfg.GetString("REDIS_PASSWORD"),
			UserTokensTTL: jsonCfg.GetDuration("repository.userTokensTTL"),
		},
		Postgres: &PostgresConfig{
			Host:     envCfg.GetString("POSTGRES_HOST"),
			User:     envCfg.GetString("POSTGRES_USER"),
			Password: envCfg.GetString("POSTGRES_PASSWORD"),
			DBName:   envCfg.GetString("POSTGRES_DB"),
			Port:     envCfg.GetInt("POSTGRES_PORT"),
		},
		TokenManager: &TokenManagerConfig{
			SigningKey: envCfg.GetString("JWT_SIGNING_KEY"),
		},
		Dune: &DuneConfig{
			ApiKey: envCfg.GetString("DUNE_API_KEY"),
		},
	}, nil
}

func (p *PostgresConfig) PgSource() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable pool_max_conns=32",
		p.Host, p.Port, p.User, p.Password, p.DBName)
}
