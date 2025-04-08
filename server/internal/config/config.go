package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	Clients    ClientsConfig  `yaml:"clients"`
	AppSecret  string         `yaml:"app_secret" env-required:"true" env:"APP_SECRET"`
	Redis      RedisConfig    `yaml:"redis"`
	Postgres   PostgresConfig `yaml:"postgres" env-required:"true"`
}

type HTTPServer struct {
	Address         string        `yaml:"address" env-default:"localhost:8081"`
	Timeout         time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"60s"`
	ShutDownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
}

type Client struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retriesCount"`
	// Insecure     bool          `yaml:"insecure"` // TODO: implement insecure in prod
}

type ClientsConfig struct {
	SSO Client `yaml:"sso"`
}

type RedisConfig struct {
	Addr            string        `yaml:"addr"`
	Password        string        `yaml:"password"`
	DB              int           `yaml:"db"`
	SSOCacheTTL     time.Duration `yaml:"sso_cache_ttl" env-default:"5m"`
	ProfileCacheTTL time.Duration `yaml:"profile_cache_ttl" env-default:"10m"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn" env:"POSTGRES_DSN" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH path is empty")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
