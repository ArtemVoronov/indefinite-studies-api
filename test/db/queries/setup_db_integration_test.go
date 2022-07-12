//go:build integration
// +build integration

package queries_test

import (
	"os"
	"testing"

	integrationTesting "github.com/ArtemVoronov/indefinite-studies-api/internal/app/testing"
)

var DuplicateKeyConstraintViolationError string = "pq: duplicate key value violates unique constraint \"%s\""

func TestMain(m *testing.M) {
	integrationTesting.Setup()
	code := m.Run()
	integrationTesting.Shutdown()
	os.Exit(code)
}
