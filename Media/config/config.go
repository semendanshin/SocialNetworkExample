package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

// Config is the configuration for the application.
type Config struct {
	Env    string `yaml:"env" env-required:"true"` // dev, test, prod
	Server Server `yaml:"server"`
	//Postgres Postgres `yaml:"postgres"`
	Minio  Minio  `yaml:"minio"`
	Tokens Tokens `yaml:"tokens"`
}

// Server is the configuration for the server.
type Server struct {
	Address string        `yaml:"address" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

// Minio is the configuration for the Minio storage.
type Minio struct {
	Endpoint  string `yaml:"endpoint" env-required:"true"`
	AccessKey string `yaml:"access_key" env-required:"true"`
	SecretKey string `yaml:"secret_key" env-required:"true"`
	UseSSL    *bool  `yaml:"use_ssl" env-required:"true"`
	Bucket    string `yaml:"bucket" env-required:"true"`
}

// Postgres is the configuration for the PostgreSQL database.
//type Postgres struct {
//	Host string `yaml:"host"`
//	Port int    `yaml:"port"`
//	User string `yaml:"user"`
//	Pass string `yaml:"pass"`
//	Name string `yaml:"name"`
//}

// Tokens is the configuration for the JWT tokens.
type Tokens struct {
	PublicKeyPath string `yaml:"public_key_path" env-required:"true"`
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
