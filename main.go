package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ArtemVoronov/indefinite-studies-api/api/rest/v1/ping"
	"github.com/ArtemVoronov/indefinite-studies-api/api/rest/v1/tasks"
	"github.com/ArtemVoronov/indefinite-studies-api/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Cors() gin.HandlerFunc {
	// TODO: for release add appropiate domains
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

// TODO: add authorization
// TODO: add graceful shutdown

func main() {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	port, portExists := os.LookupEnv("APP_PORT")
	if !portExists {
		port = "3000"
	}

	host := ":" + port

	log.Println("starting indefinite-studies-api at " + host + " ...")

	router := gin.Default()

	router.Use(Cors())

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	db.DB = db.Setup()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", ping.Ping)
		v1.GET("/tasks/", tasks.GetTasks)
		v1.GET("/tasks/:id", tasks.GetTask)
		v1.POST("/tasks/", tasks.CreateTask)
		v1.PUT("/tasks/", tasks.UpdateTask)
		v1.DELETE("/tasks/:id", tasks.DeleteTask)
	}

	defer db.DB.Close()
	router.Run(host)
}
