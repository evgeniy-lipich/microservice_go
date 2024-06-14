package env

import (
	"github.com/evgeniy-lipich/microservice_go/week_2/config/internal/config"
	"github.com/pkg/errors"
	"net"
	"os"
)

var _ config.GRPCConfig = (*grpcConfig)(nil)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

// структура имплементирует интерфейс
type grpcConfig struct {
	host string
	port string
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func NewGRPCConfig() (*grpcConfig, error) {
	// получить хост из системных переменных
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	// получить хост из системных переменных
	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	// вернуть структуру конфига
	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}
