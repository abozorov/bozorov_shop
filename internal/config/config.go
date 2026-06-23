package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerConfig
	DBConfig
}

type ServerConfig struct {
	// server
	HttpHost        string `env:"PORT"`
	Storage         string `env:"STORAGE"`
	AuditLogStorage string `env:"AUDITLOGSTORAGE"`
}

type DBConfig struct {
	// database
	DBHost     string `env:"DBHOST"`
	DBPort     string `env:"DBPORT"`
	DBUser     string `env:"DBUSER"`
	DBPassword string `env:"DBPASSWORD"`
	DBName     string `env:"DBNAME"`
}

func NewConfig(path string) (Config, error) {
	var cnf Config

	err := cleanenv.ReadConfig(path, &cnf)
	if err != nil {
		return Config{}, fmt.Errorf("cleanenv.ReadConfig: %w", err)
	}

	return cnf, nil
}
