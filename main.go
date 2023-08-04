package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

// Configuration
const dbFile = "./data/cardata.db"
const IP = "localhost"
const port = "8000"

type Car struct {
	ID        int
	Make      string
	Model     string
	BuildDate time.Time
	ColourID  int
}

type Colour struct {
	ID   int
	Name string
}

func main() {
	app := initApp() // splitting the app initialisation and listening out allows us to use app.Test for testing
	app.Listen(IP + ":" + port)
}

func initApp() *fiber.App {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Failed to connect to database: %v", dbFile)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Get("/cars", func(c *fiber.Ctx) error {
		return getCarsHandler(c, db)
	})

	app.Post("/cars", func(c *fiber.Ctx) error {
		return postCarsHandler(c, db)
	})

	app.Get("/car/:id", func(c *fiber.Ctx) error {
		return getCarByIdHandler(c, db)
	})

	app.Delete("/car/:id", func(c *fiber.Ctx) error {
		return deleteCarByIdHandler(c, db)
	})

	return app
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Hello")
}

func getCarsHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("GET to /cars")
}

func postCarsHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("POST to /cars")
}

func getCarByIdHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("GET to /car/:id")
}

func deleteCarByIdHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("DELETE to /car/:id")
}
