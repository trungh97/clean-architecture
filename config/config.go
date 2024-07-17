package config

import (
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Server   Server
		Database Database
		JWT      JWT
		Redis    Redis
		Session  Session
		Cookie   Cookie
		Logger   Logger
	}

	Server struct {
		Port         string
		Mode         string
		CSRF         bool
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		PprofPort    string
		AppVersion   string
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

	Redis struct {
		Host     string
		Port     int
		Password string
		Database int
	}

	Session struct {
		Prefix string
		Name   string
		Expire int
	}

	Cookie struct {
		Name     string
		MaxAge   int
		Secure   bool
		HTTPOnly bool
	}

	Logger struct {
		Development       bool
		DisableCaller     bool
		DisableStacktrace bool
		Encoding          string
		Level             string
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
