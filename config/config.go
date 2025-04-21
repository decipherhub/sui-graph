package config

import (
	"fmt"
	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/decipherhub/sui-graph/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	Logger *log.Entry
	Server ServerConfig `mapstructure:"server"`
	Sui    SuiConfig    `mapstructure:"sui"`

	SuiClient sui.ISuiAPI
	Database  DatabaseConfig `mapstructure:"database"`
	Graph     GraphConfig    `mapstructure:"graph"`
	Logging   LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type SuiConfig struct {
	RPCUrl           string `json:"rpc_url"`
	FetchIntervalSec int    `mapstructure:"fetch_interval_sec"`
	CheckpointStart  int64  `mapstructure:"checkpoint_start"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

type GraphConfig struct {
	MaxConcurrencyLayer        int  `mapstructure:"max_concurrency_layer"`
	EnableSharedObjectTracking bool `mapstructure:"enable_shared_object_tracking"`
}

type LoggingConfig struct {
	Level string `mapstructure:"level"`
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

	// Logger
	logger := log.New().WithFields(log.Fields{
		"version": version.AppVersion,
	})
	logger.Logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logger.Logger.SetOutput(os.Stdout)

	// SuiClient
	suiClient := sui.NewSuiClient(constant.SuiTestnetEndpoint)

	cfg := &Config{
		Logger: logger,
		Server: ServerConfig{
			Port: v.GetInt("SERVER_PORT"),
			Host: v.GetString("SERVER_HOST"),
		},
		Sui: SuiConfig{
			RPCUrl:           v.GetString("SUI_RPC_URL"),
			FetchIntervalSec: v.GetInt("FETCH_INTERVAL_SEC"),
			CheckpointStart:  v.GetInt64("CHECKPOINT_START"),
		},
		SuiClient: suiClient,
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
