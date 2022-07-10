//go:build integration
// +build integration

package queries_test

import (
	integrationTesting "github.com/ArtemVoronov/indefinite-studies-api/internal/app/testing"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	integrationTesting.Setup()
	code := m.Run()
	integrationTesting.Shutdown()
	os.Exit(code)
}
