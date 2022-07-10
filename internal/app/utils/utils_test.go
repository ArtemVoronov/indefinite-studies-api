//go:build unit
// +build unit

package utils_test

import (
	"testing"

	"fmt"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/app/utils"
	"github.com/stretchr/testify/assert"
)

func TestSha512Hash(t *testing.T) {
	//expected output of 'openssl dgst -sha512 <<< any_password_here'
	expected := "16d99c6502225c7e8ee5c85af1070cbcf04724763836ad29edaedab552a54c63d79fb04f62e7a8b4a4b849a6edc558010a67b9b57a949aaf425c6a0dc821fa2d"

	fmt.Printf("expected: %s", expected)
	actual := utils.CreateSHA512HashHexEncoded("any_password_here\n")

	assert.Equal(t, expected, actual)
}

func TestSha256Hash(t *testing.T) {
	//expected output of 'openssl dgst -sha256 <<< any_password_here'
	expected := "f5d3eb241f80e6ce4581f07ef7a770e59f07d6cbdf55b35cf1416aff2acba064"

	fmt.Printf("expected: %s", expected)
	actual := utils.CreateSHA256HashHexEncoded("any_password_here\n")

	assert.Equal(t, expected, actual)
}
