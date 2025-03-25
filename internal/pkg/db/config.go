package db

import "time"

type Config struct {
	DB             string        `env:"CLICKHOUSE_DB"`
	User           string        `env:"CLICKHOUSE_USER"`
	Password       string        `env:"CLICKHOUSE_PASSWORD"`
	Host           string        `env:"CLICKHOUSE_HOST"`
	Port           uint16        `env:"CLICKHOUSE_PORT"`
	ConnectTimeout time.Duration `yaml:"connectTimeout" env-default:"5m"`
}
