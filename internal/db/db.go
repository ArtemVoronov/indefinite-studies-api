package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// TODO: reorder funcs, types, vars in 'db' namespace

type Singleton interface {
	GetDB() *sql.DB
	GetTimeout() time.Duration
}

type singleton struct {
	rwmutex sync.RWMutex
	db      *sql.DB
	timeout time.Duration
}

var once sync.Once
var instance *singleton

func GetInstance() Singleton {
	once.Do(func() {
		if instance == nil {
			instance = new(singleton)
			instance.setDB(createDatabase())
			instance.setTimeout(30 * time.Second) // TODO: make timeout confugrable for tests
		}
	})
	return instance
}

func (s *singleton) setDB(sqldb *sql.DB) {
	s.rwmutex.Lock()
	defer s.rwmutex.Unlock()
	s.db = sqldb
}

func (s *singleton) setTimeout(timeout time.Duration) {
	s.rwmutex.Lock()
	defer s.rwmutex.Unlock()
	s.timeout = timeout
}

func (s *singleton) GetDB() *sql.DB {
	s.rwmutex.RLock()
	defer s.rwmutex.RUnlock()
	return s.db
}

func (s *singleton) GetTimeout() time.Duration {
	s.rwmutex.RLock()
	defer s.rwmutex.RUnlock()
	return s.timeout
}

var ErrorTaskDuplicateKey = errors.New("pq: duplicate key value violates unique constraint \"tasks_name_state_unique\"")
var ErrorTagDuplicateKey = errors.New("pq: duplicate key value violates unique constraint \"tags_name_state_unique\"")
var ErrorUserDuplicateKey = errors.New("pq: duplicate key value violates unique constraint \"users_email_state_unique\"")

func createDatabase() *sql.DB {

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

type SqlQueryFunc func(database *sql.DB, ctx context.Context)

func RunWithWithTimeout(f SqlQueryFunc) func() {
	database := GetInstance().GetDB()
	timeout := GetInstance().GetTimeout()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return func() {
		defer cancel()
		f(database, ctx)
	}
}
