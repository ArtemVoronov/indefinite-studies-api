//go:build integration
// +build integration

package queries_test

import (
	"fmt"
	"testing"
)

func TestUnit(t *testing.T) {
	fmt.Println("testing!!!!!!:", t.Name())
}

func TestIntegration(t *testing.T) {
	fmt.Println("testing!!!!!!:", t.Name())
}
