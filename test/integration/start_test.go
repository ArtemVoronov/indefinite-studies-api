//go:build integration
// +build integration

package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/ping"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/tasks"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var Router *gin.Engine

func TestMain(m *testing.M) {
	Setup()
	Router = SetupRouter()
	code := m.Run()
	Shutdown()
	os.Exit(code)
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", ping.Ping)

	r.GET("/tasks", tasks.GetTasks)
	r.GET("/tasks/:id", tasks.GetTask)
	r.POST("/tasks", tasks.CreateTask)
	r.PUT("/tasks/:id", tasks.UpdateTask)
	r.DELETE("/tasks/:id", tasks.DeleteTask)

	return r
}

func GetRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	d1 := path.Join(path.Dir(b))
	return d1[:len(d1)-len("/test/integration")]
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
	db.GetInstance()
}

func Shutdown() {
	defer db.GetInstance().GetDB().Close()
}