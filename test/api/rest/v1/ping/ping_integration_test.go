//go:build integration
// +build integration

package ping_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/ping"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", ping.Ping)
	return r
}

func TestApiPing(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "\"Pong!\"", w.Body.String())
}
