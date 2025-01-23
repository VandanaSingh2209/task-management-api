package Database

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func TaskManagementDb() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tskMgmtDb := os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") +
		"@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" +
		os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err = gorm.Open(mysql.Open(tskMgmtDb), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	log.Println("Database connection established!")
}

func ResetDBPoolConnection(database string, errorMsg string) bool {
	if strings.Contains(errorMsg, "sql: database is closed") || strings.Contains(errorMsg, "Server shutdown in progress") || strings.Contains(errorMsg, "invalid connection") {

		switch database {
		case "taskManagementDB":
			TaskManagementDb()
		}

		return true
	}
	return false
}
