package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/ping"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/tasks"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	initEnv()
	apiUsers := getApiUsers()
	host := getHost()

	router := gin.Default()

	router.Use(cors())

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

	v1 := router.Group("/api/v1", gin.BasicAuth(apiUsers))
	{
		v1.GET("/ping", ping.Ping)
		v1.GET("/tasks/", tasks.GetTasks)
		v1.GET("/tasks/:id", tasks.GetTask)
		v1.POST("/tasks/", tasks.CreateTask)
		v1.PUT("/tasks/", tasks.UpdateTask)
		v1.DELETE("/tasks/:id", tasks.DeleteTask)
	}

	startServer(host, router)
}

func cors() gin.HandlerFunc {
	// TODO: for release add appropiate domains
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func initEnv() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func getHost() string {
	port, portExists := os.LookupEnv("APP_PORT")
	if !portExists {
		port = "3000"
	}
	host := ":" + port
	return host
}

func getApiUsers() gin.Accounts {
	apiKey, apiKeyExists := os.LookupEnv("AUTH_USERNAME")
	if !apiKeyExists {
		log.Fatalf("Missed enviroment variable: %s. Check the .env file or OS enviroment vars", "AUTH_USERNAME")
	}

	apiAuthUser, apiAuthUserExists := os.LookupEnv("AUTH_PASSWORD")
	if !apiAuthUserExists {
		log.Fatalf("Missed enviroment variable: %s. Check the .env file or OS enviroment vars", "AUTH_PASSWORD")
	}
	apiUsers := gin.Accounts{apiKey: apiAuthUser}
	return apiUsers
}

func startServer(host string, router *gin.Engine) {
	srv := &http.Server{
		Addr:    host,
		Handler: router,
	}

	// Initializing the server in a goroutine so that it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	defer db.DB.Close()
	log.Println("Server exiting")
}
