package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

const (
	StorageModePostgres = "postgres"
	StorageModeMemory   = "memory"
)

type ServerConfig struct {
	Logger   LoggerConf
	Endpoint EndpointConf
	Storage  StorageConf
}

type LoggerConf struct {
	Level  string
	Output string
}

type EndpointConf struct {
	Host     string
	HTTPPort int
	GRPCPort int
}

type StorageConf struct {
	Mode     string
	Host     string
	Port     int
	DBName   string
	User     string
	Password string
}

func NewServerConfig(filePath string) (c ServerConfig) {
	_, err := toml.DecodeFile(filePath, &c)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	return c
}
