package database

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"instagram-clone.com/m/config"
)

type mysqlDatabase struct {
	Db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *mysqlDatabase
)

func NewMySQLDatabase(conf *config.Config) *mysqlDatabase {
	once.Do(func() {
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			conf.Database.User,
			conf.Database.Password,
			conf.Database.Host,
			conf.Database.Port,
			conf.Database.Database,
		)

		db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		dbInstance = &mysqlDatabase{
			Db: db,
		}

	})

	return dbInstance
}

func (p *mysqlDatabase) GetDb() *gorm.DB {
	return dbInstance.Db
}
