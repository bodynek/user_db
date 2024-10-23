package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadConfig() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	viper.AutomaticEnv() // Automatically override config with environment variables
}

func ConnectDB() (*gorm.DB, error) {
	LoadConfig()

	host := viper.GetString("DB_HOST")
	user := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	dbname := viper.GetString("DB_NAME")
	port := viper.GetString("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return db, nil
}

func GetServerPort() string {
	LoadConfig()
	port := viper.GetString("APP_PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}
	return port
}
