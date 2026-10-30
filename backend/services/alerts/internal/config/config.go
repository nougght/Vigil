package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/nougght/monitoring-system/backend/pkg/storage"
	"github.com/nougght/monitoring-system/pkg/util"
)

type HTTPConfig struct {
	ServerPort int
}
type GRPCConfig struct {
	ServerPort int
}

type Config struct {
	Postgres *storage.PostgresConfig
	HTTP     *HTTPConfig
	//GRPC     *GRPCConfig
	NatsURL string
}

func MustLoadConfig(path string) *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(err.Error())
	}
	cfg := &Config{}

	cfg.Postgres = &storage.PostgresConfig{
		Host:     util.MustGetEnvVar("POSTGRES_HOST"),
		Port:     util.MustGetEnvVar("POSTGRES_PORT"),
		User:     util.MustGetEnvVar("POSTGRES_USER"),
		Password: util.MustGetEnvVar("POSTGRES_PASSWORD"),
		DBName:   util.MustGetEnvVar("POSTGRES_DB"),
		SSLMode:  util.MustGetEnvVar("POSTGRES_SSL_MODE"),
	}

	serverPort := util.MustGetPortEnvVar("HTTP_SERVER_PORT")
	cfg.HTTP = &HTTPConfig{
		ServerPort: serverPort,
	}

	//grpcPort := util.MustGetPortEnvVar("GRPC_SERVER_PORT")
	//cfg.GRPC = &GRPCConfig{
	//	ServerPort: grpcPort,
	//}

	natsURL := util.MustGetEnvVar("NATS_URL")
	cfg.NatsURL = natsURL

	return cfg
}
