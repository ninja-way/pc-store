package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Service ServiceSettings

	DB DBSettings

	Environment string        `mapstructure:"environment"`
	CacheTTL    time.Duration `mapstructure:"cache_ttl"`

	Auth Auth `mapstructure:"auth"`

	Server Server `mapstructure:"server"`
}

type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServiceSettings struct {
	HashSalt    string
	TokenSecret string
}

type Auth struct {
	TokenTTL time.Duration `mapstructure:"token_ttl"`
}

type DBSettings struct {
	Host     string
	Port     int
	UserName string
	Password string
	SSLMode  string
	DBName   string
}

func New(folder, filename string) (*Config, error) {
	cfg := new(Config)

	viper.AddConfigPath(folder)
	viper.SetConfigName(filename)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if err := envconfig.Process("pc_service", &cfg.Service); err != nil {
		return nil, err
	}

	if err := envconfig.Process("pc", &cfg.DB); err != nil {
		return nil, err
	}

	return cfg, nil
}
