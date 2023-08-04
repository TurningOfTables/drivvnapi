package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

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
	app.Listen("localhost:8000")
}

func initApp() *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c)
	})

	app.Get("/cars", func(c *fiber.Ctx) error {
		return getCarsHandler(c)
	})

	app.Post("/cars", func(c *fiber.Ctx) error {
		return postCarsHandler(c)
	})

	app.Get("/car/:id", func(c *fiber.Ctx) error {
		return getCarByIdHandler(c)
	})

	app.Delete("/car/:id", func(c *fiber.Ctx) error {
		return deleteCarByIdHandler(c)
	})

	return app
}

func indexHandler(c *fiber.Ctx) error {
	return c.SendString("Hello")
}

func getCarsHandler(c *fiber.Ctx) error {
	return c.SendString("GET to /cars")
}

func postCarsHandler(c *fiber.Ctx) error {
	return c.SendString("POST to /cars")
}

func getCarByIdHandler(c *fiber.Ctx) error {
	return c.SendString("GET to /car/:id")
}

func deleteCarByIdHandler(c *fiber.Ctx) error {
	return c.SendString("DELETE to /car/:id")
}
