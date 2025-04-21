package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Server   ServerConfig
	Sui      SuiConfig
	Database DatabaseConfig
	Graph    GraphConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port int
	Host string
}

type SuiConfig struct {
	RPCUrl           string
	FetchIntervalSec int
	CheckpointStart  int64
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type GraphConfig struct {
	MaxConcurrencyLayer        int
	EnableSharedObjectTracking bool
}

type LoggingConfig struct {
	Level string
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Load from .env
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(configPath)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := v.ReadInConfig() // Read .env
	if err != nil {
		fmt.Println("No .env file found or error:", err)
	}

	// Optional: Also try to load config.yaml
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	if err := v.MergeInConfig(); err != nil {
		fmt.Println("No config.yaml found or error:", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: v.GetInt("SERVER_PORT"),
			Host: v.GetString("SERVER_HOST"),
		},
		Sui: SuiConfig{
			RPCUrl:           v.GetString("SUI_RPC_URL"),
			FetchIntervalSec: v.GetInt("FETCH_INTERVAL_SEC"),
			CheckpointStart:  v.GetInt64("CHECKPOINT_START"),
		},
		Database: DatabaseConfig{
			Driver: v.GetString("DB_DRIVER"),
			DSN:    v.GetString("DB_DSN"),
		},
		Graph: GraphConfig{
			MaxConcurrencyLayer:        v.GetInt("MAX_CONCURRENCY_LAYER"),
			EnableSharedObjectTracking: v.GetBool("ENABLE_SHARED_OBJECT_TRACKING"),
		},
		Logging: LoggingConfig{
			Level: v.GetString("LOG_LEVEL"),
		},
	}

	return cfg, nil
}
