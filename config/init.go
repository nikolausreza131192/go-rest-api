package config

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Config define app config model
type Config struct {
	Main     Main
	Database map[string]Database
}

// Main define app main configuration
type Main struct {
	APIPort string
}

// Database define database model
type Database struct {
	Driver   string
	User     string
	Password string
	Name     string
}

// InitConfig Initialize app configuration
func InitConfig() Config {
	fmt.Println("Init config...")

	//Load environmenatal variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error reading .env file")
	}

	return Config{
		Main: Main{
			APIPort: os.Getenv("APP_PORT"),
		},
		Database: map[string]Database{
			"stone_work": Database{
				Driver:   "mysql",
				Name:     os.Getenv("DB_NAME"),
				User:     os.Getenv("DB_USERNAME"),
				Password: os.Getenv("DB_PASSWORD"),
			},
		},
	}
}
