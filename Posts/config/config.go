package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

// Config is the configuration for the application.
type Config struct {
	Env         string   `yaml:"env" env-required:"true"` // dev, test, prod
	Server      Server   `yaml:"server"`
	UseDatabase *bool    `yaml:"use_database" env-required:"false" env-default:"false"`
	Postgres    Postgres `yaml:"postgres"`
	Tokens      Tokens   `yaml:"tokens"`
}

// Server is the configuration for the server.
type Server struct {
	Address string        `yaml:"address" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

// Postgres is the configuration for the PostgreSQL database.
type Postgres struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Name string `yaml:"name"`
}

// Tokens is the configuration for the JWT tokens.
type Tokens struct {
	Secret     string        `yaml:"secret" env-required:"true"`
	AccessTTL  time.Duration `yaml:"access_ttl" env-required:"true"`
	RefreshTTL time.Duration `yaml:"refresh_ttl" env-required:"true"`
}

// MustParseConfig parses the configuration from the given path.
func MustParseConfig(path string) Config {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

// FetchPath fetches the path to the configuration file.
func FetchPath() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	if path == "" {
		path = "config/local.yaml"
	}

	return path
}
