package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB // TODO: make injection of DB to api functions

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