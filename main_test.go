package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

	assert.Equal(t, 200, resp.StatusCode)
}

// Test GET to "/cars"
func TestGetCars(t *testing.T) {
	app := initApp()
	req := httptest.NewRequest(http.MethodGet, "/cars", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 200, resp.StatusCode)
}

// Test POST to "/cars"
func TestPostCars(t *testing.T) {
	postBody := []Car{{Make: "TestMake", Model: "TestModel", BuildDate: time.Now().Format(time.DateOnly), ColourID: 1}}
	json, err := json.Marshal(postBody)
	if err != nil {
		t.Error("Error encoding JSON body")
	}
	reader := bytes.NewReader(json)

	app := initApp()
	req := httptest.NewRequest(http.MethodPost, "/cars", reader)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 201, resp.StatusCode)
}

// Test GET to "/car/:id"
func TestGetCarById(t *testing.T) {
	app := initApp()
	req := httptest.NewRequest(http.MethodGet, "/car/1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 200, resp.StatusCode)
}

// Test DELETE to "/car/:id"
func TestDeleteCarById(t *testing.T) {
	app := initApp()
	req := httptest.NewRequest(http.MethodDelete, "/car/1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 200, resp.StatusCode)
}

func TestCarValidation(t *testing.T) {
	var c = Car{
		Make:      "BMW",
		Model:     "3 Series",
		BuildDate: "2020-01-20",
		ColourID:  2,
	}

	err := carValidation(c)
	if err != nil {
		t.Error(err)
	}
}

func TestFailedCarValidation(t *testing.T) {

	// Deliberately missing a required field to force a validation error
	var c = Car{
		Make:      "BMW",
		BuildDate: "2020-01-20",
		ColourID:  2,
	}

	err := carValidation(c)
	assert.NotNil(t, err)
}

func TestBuildDateValidation(t *testing.T) {
	err := buildDateValidation("2020-01-20")
	assert.Nil(t, err)
}

func TestBuildDateValidationTooOld(t *testing.T) {
	err := buildDateValidation("2018-01-20")
	assert.NotNil(t, err)
}

func respBodyToString(r io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.String()
}
