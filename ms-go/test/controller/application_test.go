package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ms-go/router"

	"github.com/gin-gonic/gin"
	"github.com/likexian/gokit/assert"
)

func TestIndexHome(t *testing.T) {
	// Build our expected body
	expectedBody := gin.H{
		"message": "[ms-go] | Success",
		"status":  200,
	}

	// Grab our router
	r := router.SetupRouter()

	// Perform a GET request
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Decode response body
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	// Assert message and status in response body
	assert.Equal(t, expectedBody["message"], responseBody["message"])
	assert.Equal(t, expectedBody["status"], responseBody["status"])
}

func TestNotFound(t *testing.T) {
	// Grab our router
	r := router.SetupRouter()

	// Perform a GET request to a non-existent route
	req, _ := http.NewRequest(http.MethodGet, "/non-existent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert status code
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostNotAllowed(t *testing.T) {
	// Grab our router
	r := router.SetupRouter()

	// Perform a POST request
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert status code
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
