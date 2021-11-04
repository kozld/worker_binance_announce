package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type WorkerConfig struct {
}

type TraderConfig struct {
	QuantityUSDT    float64 `env:"QUANTITY_USDT,required"`
	GateIOApiKey    string  `env:"GATEIO_API_KEY,required"`
	GateIOApiSecret string  `env:"GATEIO_API_SECRET,required"`
}

type DatabaseConfig struct {
	PostgresHost     string `env:"POSTGRES_HOST" envDefault:"postgres"`
	PostgresPort     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresUser     string `env:"POSTGRES_USER" envDefault:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresDbName   string `env:"POSTGRES_DB" envDefault:"binance_api"`
}

func GetWorkerConfig() *WorkerConfig {
	cfg := &WorkerConfig{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("Cannot parse initial ENV vars: ", err)
	}

	return cfg
}

func GetTraderConfig() *TraderConfig {
	cfg := &TraderConfig{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("Cannot parse initial ENV vars: ", err)
	}

	return cfg
}

func GetDatabaseConfig() *DatabaseConfig {
	cfg := &DatabaseConfig{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("Cannot parse initial ENV vars: ", err)
	}

	return cfg
}
