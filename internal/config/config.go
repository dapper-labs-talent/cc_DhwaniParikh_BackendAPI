package config

import (
	"github.com/rs/zerolog"
	"strings"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/rs/zerolog/log"
)

type ServerConfig struct {
	Port int `koanf:"Port" validate:"required"`
}

type DatabaseConfig struct {
	ConnectionString string `koanf:"ConnectionString" validate:"required"`
	LogSQL           bool   `koanf:"LogSQL"`
}

type AuthConfig struct {
	JWT JWTConfig `koanf:"JWT" validate:"required"`
}

type JWTConfig struct {
	Key string `koanf:"Key" validate:"required"`
}

type FunctionConfig struct {
	APIKey    string `koanf:"APIKey" validate:"required"`
	APISecret string `koanf:"APISecret" validate:"required"`
}

type Config struct {
	Server   ServerConfig   `koanf:"Server"`
	Database DatabaseConfig `koanf:"Database"`
	Auth     AuthConfig     `koanf:"Auth"`
	Log      LogConfig      `koanf:"Log"`
}
type LogConfig struct {
	Level        string `koanf:"Level"`
	ZeroLogLevel zerolog.Level
}

// Prefix for environment variables
var envPrefix = "DAPPER_"

func LoadConfig() (*Config, error) {
	config := &Config{}

	_ = godotenv.Overload()

	k := koanf.New(".")

	log.Info().Msg("Loading configuration from environment")
	// Load environment variables and replace _ with . to conform to the keys
	err := k.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.TrimPrefix(s, envPrefix), "_", ".")
	}), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal config from environment")
	}

	log.Info().Msg("Unmarshalling configuration from environment")
	err = k.UnmarshalWithConf("", &config, koanf.UnmarshalConf{Tag: "koanf"})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal config from environment")
	}

	log.Info().Msg("Validating configuration")
	return config, nil
}
