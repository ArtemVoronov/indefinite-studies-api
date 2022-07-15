package db_test

// TODO: finish tests for initiation of db service
// TODO: add mocks for *sql.DB

// import (
// 	"testing"

// 	integrationTesting "github.com/ArtemVoronov/indefinite-studies-api/internal/app/testing"
// 	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
// 	"github.com/stretchr/testify/assert"
// )

// func TestGetInstance(t *testing.T) {
// 	integrationTesting.InitTestEnv()
// 	singleton1 := db.GetInstance()

// 	assert.NotNil(t, singleton1)
// 	assert.NotNil(t, singleton1.GetDB())

// 	singleton2 := db.GetInstance()

// 	assert.NotNil(t, singleton2)
// 	assert.NotNil(t, singleton2.GetDB())
// 	assert.Equal(t, singleton1, singleton2)
// 	assert.Equal(t, singleton1.GetDB(), singleton2.GetDB())
// }

// func TestParallel(t *testing.T) {
// 	integrationTesting.InitTestEnv()

// 	var singleton1 db.Singleton
// 	var singleton2 db.Singleton

// 	var wg sync.WaitGroup
// 	for i := 0; i < 5000; i++ {
// 		wg.Add(1)
// 		go func() {
// 			if singleton1 != nil && singleton2 != nil {
// 				assert.NotNil(t, singleton1)
// 				assert.NotNil(t, singleton1.GetDB())
// 				assert.NotNil(t, singleton2)
// 				assert.NotNil(t, singleton2.GetDB())
// 				assert.Equal(t, singleton1, singleton2)
// 				assert.Equal(t, singleton1.GetDB(), singleton2.GetDB())
// 			}
// 			singleton1 = db.GetInstance()

// 			wg.Done()
// 		}()
// 		wg.Add(1)
// 		go func() {
// 			if singleton1 != nil && singleton2 != nil {
// 				assert.NotNil(t, singleton1)
// 				assert.NotNil(t, singleton1.GetDB())
// 				assert.NotNil(t, singleton2)
// 				assert.NotNil(t, singleton2.GetDB())
// 				assert.Equal(t, singleton1, singleton2)
// 				assert.Equal(t, singleton1.GetDB(), singleton2.GetDB())
// 			}

// 			singleton2 = db.GetInstance()

// 			wg.Done()
// 		}()
// 	}
// 	wg.Wait()
// }
