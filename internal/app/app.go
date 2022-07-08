package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
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

func InitEnv() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func GetHost() string {
	port, portExists := os.LookupEnv("APP_PORT")
	if !portExists {
		port = "3000"
	}
	host := ":" + port
	return host
}

func GetApiUsers() gin.Accounts {
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

func StartServer(host string, router *gin.Engine) {
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
