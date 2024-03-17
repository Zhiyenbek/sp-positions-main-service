package config

import (
	"time"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

type Configs struct {
	App   *AppConfig `json:"app" mapstructure:"app"`
	DB    *DBConf    `json:"db" mapstructure:"db"`
	Redis *RedisConf `json:"redis" mapstructure:"redis"`
	Token *Token     `json:"token" mapstructure:"token"`
}

type AppConfig struct {
	TimeOut time.Duration `json:"timeout" mapstructure:"timeout"`
	Port    int           `json:"port" mapstructure:"port"`
}

type DBConf struct {
	Host     string        `json:"host" mapstructure:"host"`
	Port     int           `json:"port" mapstructure:"port"`
	Username string        `json:"username" mapstructure:"user"`
	Password string        `json:"password" mapstructure:"password"`
	DBName   string        `json:"dbname" mapstructure:"db_name"`
	SSLMode  string        `json:"sslmode" mapstructure:"ssl_mode"`
	TimeOut  time.Duration `json:"timeout" mapstructure:"timeout"`
}

type RedisConf struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     int    `json:"port" mapstructure:"port"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	DB       int    `json:"db" mapstructure:"db"`
}

type Token struct {
	Refresh *TokenConf `json:"refresh" mapstructure:"refresh"`
	Access  *TokenConf `json:"access" mapstructure:"access"`
}

type TokenConf struct {
	TokenSecret string        `json:"token_secret" mapstructure:"token_secret"`
	TTL         time.Duration `json:"ttl"          mapstructure:"ttl"`
	Domain      string        `json:"domain"       mapstructure:"domain"`
}

func New() (*Configs, error) {
	configFile := "../config/config.yaml"
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Configs{}

	if err := defaults.Set(cfg); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
