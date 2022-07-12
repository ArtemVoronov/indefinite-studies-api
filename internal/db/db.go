package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// TODO: add explicit tx using, and context with 30 second timeout, make timeout and db *sql.D as one struct that is used by queries
// TODO: then make injection of DB to api functions and integration tests, maybe just a singleton
var DB *sql.DB

var ErrorDuplicateKey = errors.New("pq: duplicate key value violates unique constraint \"tasks_name_state_unique\"")

func Setup() *sql.DB {

	dbEnvVars := [6]string{"DATABASE_HOST", "DATABASE_PORT", "DATABASE_USER", "DATABASE_PASSWORD", "DATABASE_NAME", "DATABASE_SSL_MODE"}
	var variables []string
	for _, element := range dbEnvVars {

		value, variableExists := os.LookupEnv(element)
		if !variableExists {
			log.Fatalf("Missed enviroment variable: %s. Check the .env file or OS enviroment vars", element)
		}
		variables = append(variables, value)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", variables[0], variables[1], variables[2], variables[3], variables[4], variables[5])
	result, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database : %s", err)
	}

	log.Printf("----- Database service setup succeed. Database name: %s -----", variables[4])

	return result
}
