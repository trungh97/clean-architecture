package main

import (
	"log"

	"instagram-clone.com/m/config"
	"instagram-clone.com/m/database"
	"instagram-clone.com/m/internal/server"
	"instagram-clone.com/m/pkg/db/redis"
	"instagram-clone.com/m/pkg/logger"
)

// @title Instagram Clone REST API
// @version 1.0
// @description Instagram Clone REST API
// @contact.name Trung Hoang
// @BasePath /api/v1
func main() {
	log.Println("Starting api server...")

	cfg := config.GetConfig()

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	mySqlDB := database.NewMySQLDatabase(cfg)

	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis connected successfully")

	s := server.NewServer(cfg, mySqlDB, redisClient, appLogger)

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
