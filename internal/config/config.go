package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env         string         `yaml:"env" env-default:"local"`
	Port        int            `yaml:"port" env-required:"true"`
	StorageType string         `yaml:"storage_type" env-default:"inmemory"`
	Storage     DatabaseConfig `yaml:"storage" env-required:"true"`
}

type DatabaseConfig struct {
	Address  string `yaml:"db_address" env-required:"true"`
	Name     string `yaml:"db_name" env-required:"true"`
	User     string `yaml:"db_user" env-required:"true"`
	Password string `yaml:"db_password" env-required:"true"`
	SSLMode  string `yaml:"db_sslmode" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config path does not exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("falied to read config: " + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "config file path")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
