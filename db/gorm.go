package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"stockcontent-monitor-demo-back/env"
)

var Client = getDB()

func getDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(env.DatabaseConnection), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(15)
	sqlDB.SetMaxOpenConns(15)
	return db
}
