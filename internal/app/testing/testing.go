package testing

import (
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"testing"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/joho/godotenv"
)

func GetRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	d1 := path.Join(path.Dir(b))
	return d1[:len(d1)-len("/internal/app/testing")]
}

func InitTestEnv() {
	if err := godotenv.Load(GetRootPath() + "/.env.test"); err != nil {
		fmt.Println("No .env.test file found")
	}
}

func RecreateTestDB() {
	// TODO: think about carelessness removing prod database
	cmd := exec.Command("docker-compose", "--env-file", "./.env.test", "--profile", "integration-tests-only", "up", "liquibase_rollback_all_and_create_db_again")
	cmd.Dir = GetRootPath()
	_ /*stdout*/, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// uncomment for debugging
	// fmt.Println("-------------------------------------")
	// fmt.Println(string(stdout))
	// fmt.Println("-------------------------------------")
}

type TestFunc func(t *testing.T)

func RunWithRecreateDB(f TestFunc) func(t *testing.T) {
	RecreateTestDB()
	return func(t *testing.T) {
		f(t)
	}
}

func Setup() {
	InitTestEnv()
	db.DB = db.Setup()
}

func Shutdown() {
	defer db.DB.Close()
}
