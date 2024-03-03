package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"postgresHelper/lib/grpc_server"
)

const pathToConfig = "etc/config.yaml"

type Config struct {
	GRPC grpc_server.GRPCConfig `yaml:"grpc"`
	PG   Postgres               `yaml:"postgres"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

func (p *Postgres) GetConnectionInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", p.Host, p.Port, p.User, p.Password, p.Database)
}

var ConfigStruct Config

func Init() error {
	rawYaml, err := os.ReadFile(pathToConfig)
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}

	if err = yaml.Unmarshal(rawYaml, &ConfigStruct); err != nil {
		return fmt.Errorf("yaml.Unmarshal: %w", err)
	}
	return nil
}
