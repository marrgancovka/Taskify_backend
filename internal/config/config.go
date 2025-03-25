package config

import (
	"TaskTracker/internal/pkg/db"
	"TaskTracker/internal/pkg/server"
	"TaskTracker/internal/pkg/tokenizer"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
	"log"
	"os"
)

type Config struct {
	ConfigPath string `env:"CONFIG_PATH" env-default:"config/config.yaml"`

	HTTPServer server.Config    `yaml:"httpServer"`
	DB         db.Config        `yaml:"db"`
	Tokenizer  tokenizer.Config `yaml:"tokenizer"`
}

type Out struct {
	fx.Out

	HTTPServer server.Config
	DB         db.Config
	Tokenizer  tokenizer.Config
}

func MustLoad() Out {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Printf("cannot read .env file: %s", err)
		os.Exit(1)
	}

	if _, err := os.Stat(cfg.ConfigPath); os.IsNotExist(err) {
		log.Printf("config file does not exist: %s", cfg.ConfigPath)
		os.Exit(1)
	}

	if err := cleanenv.ReadConfig(cfg.ConfigPath, &cfg); err != nil {
		log.Printf("cannot read %s: %v", cfg.ConfigPath, err)
		os.Exit(1)
	}

	return Out{
		HTTPServer: cfg.HTTPServer,
		DB:         cfg.DB,
		Tokenizer:  cfg.Tokenizer,
	}
}
