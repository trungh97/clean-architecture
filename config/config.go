package config

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Server   Server
		Database Database
		JWT      JWT
	}

	Server struct {
		Port int
	}

	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}

	JWT struct {
		Secret string
	}
)

var (
	once   sync.Once
	config *Config
)

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}
	})

	return config
}
