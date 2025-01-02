package config

import (
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	err := godotenv.Load()
	if err != nil {
		return
	}

	db_conn, db_err := ConnectDB()
	if db_err != nil {
		return
	}
	DB = db_conn

}
