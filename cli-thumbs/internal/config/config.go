package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Client struct {
	Address      string        `yaml:"address" env-required:"true"`
	Timeout      time.Duration `yaml:"timeout" env-required:"true"`
	RetriesCount int           `yaml:"retries_count" env-required:"true"`
	// Insecure     bool          `yaml:"insecure" env-required:"true"`
}

type Config struct {
	Clients ClientsConfig `yaml:"clients" env-required:"true"`
}

type ClientsConfig struct {
	Thumbs Client `yaml:"thumbs" env-required:"true"`
}

func MustLoad() *Config {
	var cfg Config

	path := "./config/local.yml"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist")
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}
