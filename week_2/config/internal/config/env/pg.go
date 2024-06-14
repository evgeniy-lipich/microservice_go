package env

import (
	"github.com/evgeniy-lipich/microservice_go/week_2/config/internal/config"
	"github.com/pkg/errors"
	"os"
)

var _ config.PGConfig = (*pgConfig)(nil)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}

func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New(dsnEnvName + " is not set")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}
