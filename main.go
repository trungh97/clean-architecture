package main

import (
	"instagram-clone.com/m/config"
	"instagram-clone.com/m/database"
	"instagram-clone.com/m/server"
)

func main() {
	conf := config.GetConfig()
	db := database.NewMySQLDatabase(conf)
	server := server.NewEchoServer(conf, db)

	server.Start()
}
