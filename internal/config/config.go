package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

const (
	StorageModePostgres = "postgres"
	StorageModeMemory   = "memory"
)

type MigratorConfig struct {
	Logger  LoggerConf
	Storage StorageConf
}

type ServerConfig struct {
	Logger   LoggerConf
	Endpoint EndpointConf
	Storage  StorageConf
	Cache    CacheConf
	Queue    QueueProducerConf
}

type LoggerConf struct {
	Level  string
	Output string
}

type EndpointConf struct {
	Host     string
	GRPCPort int
}

type StorageConf struct {
	Host     string
	Port     int
	DBName   string
	User     string
	Password string
}

type CacheConf struct {
	L1Capacity int
}

type QueueServerConf struct {
	Host     string
	Port     int
	User     string
	Password string
}
type QueueProducerConf struct {
	QueueServerConf
	Exchange   string
	RoutingKey string
}

func NewRotatorConfig(filePath string) (c ServerConfig) {
	_, err := toml.DecodeFile(filePath, &c)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	return c
}

func NewMigratorConfig(filePath string) (c MigratorConfig) {
	_, err := toml.DecodeFile(filePath, &c)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	return c
}
