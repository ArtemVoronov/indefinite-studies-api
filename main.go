package main

import (
	"log"
	"os"

	"github.com/ArtemVoronov/indefinite-studies-api/api/rest/v1/ping"
	"github.com/ArtemVoronov/indefinite-studies-api/api/rest/v1/tasks"
	"github.com/ArtemVoronov/indefinite-studies-api/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

// TODO: add authorization
// TODO: update README

func main() {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	port, portExists := os.LookupEnv("APP_PORT")
	if !portExists {
		port = "3000"
	}
	host := "localhost:" + port
	log.Println("starting indefinite-studies-api at " + host + " ...")
	router := gin.Default()
	router.Use(Cors())
	db.DB = db.Setup()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", ping.Ping)
		v1.GET("/tasks/", tasks.GetTasks)
		v1.GET("/tasks/:id", tasks.GetTask)
	}

	defer db.DB.Close()
	router.Run(host)
}
