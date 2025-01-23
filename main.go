package main

import (
	"golang-assesment/Database"
	"golang-assesment/Routes"
	"io"
	"path"
	"time"

	"github.com/gin-gonic/gin"

	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	now := time.Now() //or time.Now().UTC()log
	logFileName := now.Format("2006-01-02") + ".log"
	file, err := os.OpenFile(path.Join("./storage/logs", logFileName), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Error(err)
	}
	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{
		DisableHTMLEscape: true,
		PrettyPrint:       true,
		TimestampFormat:   "2006-01-02 15:04:05",
	})
	gin.DefaultWriter = io.MultiWriter(file)
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	log.Debug("Current Time: ", now)
	log.Debug("logFileName: ", logFileName)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Error loading .env file")
	}
	//router := Routes.SetupRouter()
	router := gin.Default()

	// Set up routes
	Routes.SetupRoutes(router)
	Database.TaskManagementDb()

	// router.Static("/apidoc", "./docs")
	// router.Static("/assets", "./assets")
	router.Run(":" + os.Getenv("APP_PORT"))
}
