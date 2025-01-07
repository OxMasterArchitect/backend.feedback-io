package config

import (
	"fmt"
	"log"
	"os"

	"feedback-io.backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	User string
	Pass string
	Name string
	Host string
	Port string
}

func GetDBConfig() (*DBConfig, error) {
	config := &DBConfig{
		User: os.Getenv("MYSQL_DBUSER"),
		Pass: os.Getenv("MYSQL_DBPASSWORD"),
		Name: os.Getenv("MYSQL_DBNAME"),
		Host: os.Getenv("MYSQL_DBHOST"),
		Port: os.Getenv("MYSQL_DBPORT"),
	}

	if config.Name == "" ||
		config.Host == "" ||
		config.Pass == "" ||
		config.User == "" {
		return nil, fmt.Errorf("one or more environment variables is not sedt")
	}

	return config, nil

}

func ConnectDB() (*gorm.DB, error) {
	db_config, config_err := GetDBConfig()
	if config_err != nil {
		return nil, config_err
	}

	db_str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db_config.User, db_config.Pass, db_config.Host, db_config.Port, db_config.Name)
	db_conn, db_err := gorm.Open(mysql.Open(db_str), &gorm.Config{})

	if db_err != nil {
		return nil, db_err
	}

	log.Printf("Connected to [%s] at -> %s:%s", db_config.Name, db_config.Host, db_config.Port)
	return db_conn, nil

}

func AutoMigrateDB(DB *gorm.DB) {
	err := DB.Debug().AutoMigrate(
		&models.Suggestion{},
		&models.Comment{},
		&models.Reply{},
		&models.User{},
		&models.Category{},
		&models.Status{},
	)
	if err != nil {
		log.Fatalf("Error occured migrating database: %v", err)
	}

}
