package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test GET to "/"
func TestIndex(t *testing.T) {
	app := initApp()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	body := respBodyToString(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Hello", body)

}

// Test GET to "/cars"
func TestGetCars(t *testing.T) {
	app := initApp()
	req := httptest.NewRequest(http.MethodGet, "/cars", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	body := respBodyToString(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "GET to /cars", body)

}

// Test POST to "/cars"
func TestPostCars(t *testing.T) {
	app := initApp()
	req := httptest.NewRequest(http.MethodPost, "/cars", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	body := respBodyToString(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "POST to /cars", body)
}

// Test GET to "/car/:id"
func TestGetCarById(t *testing.T) {
	app := initApp()
	req := httptest.NewRequest(http.MethodGet, "/car/1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	body := respBodyToString(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "GET to /car/:id", body)
}

// Test POST to "/cars"
func TestDeleteCarById(t *testing.T) {
	app := initApp()
	req := httptest.NewRequest(http.MethodDelete, "/car/1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	body := respBodyToString(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "DELETE to /car/:id", body)
}

func respBodyToString(r io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.String()
}
