package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2/log"
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

	var tests = []struct {
		Date          string
		ErrorExpected bool
	}{
		{
			Date:          time.Now().Format(time.DateOnly),
			ErrorExpected: false,
		},
		{
			Date:          "1999-01-20",
			ErrorExpected: true,
		},
	}

	for _, test := range tests {
		err := buildDateValidation(test.Date)
		if test.ErrorExpected {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestColourValidation(t *testing.T) {
	ResetDB()
	db, err := ConnectToDb()
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	var tests = []struct {
		TestCar       Car
		ErrorExpected bool
	}{
		{
			TestCar:       Car{ColourID: 1},
			ErrorExpected: false,
		},
		{
			TestCar:       Car{ColourID: 9999},
			ErrorExpected: true,
		},
	}

	for _, test := range tests {
		err := colourValidation(test.TestCar, db)
		if test.ErrorExpected {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}

}
