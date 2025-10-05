package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-default:"1h"`
	Grpc        GRPCConfig    `yaml:"grpc" env-default:"{port: 8081, timeout: 10h}"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func NewConfig() *Config {
	path := fetchFlag()
	if path == "" {
		panic("config file path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found")
	}

	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("config file error" + err.Error())
	}

	return &cfg
}

func fetchFlag() string {
	var str string

	flag.StringVar(&str, "config", "", "path to config file")
	flag.Parse()

	return str
}
