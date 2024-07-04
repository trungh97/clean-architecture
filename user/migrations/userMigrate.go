package main

import (
	"instagram-clone.com/m/config"
	"instagram-clone.com/m/database"
	"instagram-clone.com/m/user/entities"
)

func main() {
	conf := config.GetConfig()
	db := database.NewMySQLDatabase(conf)

	userMigrate(db)
}

func userMigrate(db database.Database) {
	db.GetDb().AutoMigrate(&entities.User{})
}
