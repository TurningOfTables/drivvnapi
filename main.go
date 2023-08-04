package main

import (
	"database/sql"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/mattn/go-sqlite3"
)

// Configuration
const dbFile = "./data/cardata.db"
const IP = "localhost"
const port = "8000"

type Car struct {
	ID        int
	Make      string `validate:"required"`
	Model     string `validate:"required"`
	BuildDate string `validate:"required"`
	ColourID  int    `validate:"required"`
}

type Colour struct {
	ID   int    `validate:"required"`
	Name string `validate:"required"`
}

var validate *validator.Validate

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
	var cars []Car

	rows, err := db.Query("SELECT * from cars")
	if err != nil {
		log.Warn(err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error retrieving cars")
	}
	defer rows.Close()

	for rows.Next() {
		var car Car
		err := rows.Scan(&car.ID, &car.Make, &car.Model, &car.BuildDate, &car.ColourID)
		if err != nil {
			log.Warn(err)
			return c.Status(fiber.StatusInternalServerError).JSON("Error retrieving cars")
		}
		cars = append(cars, car)
	}

	return c.JSON(cars)
}

func postCarsHandler(c *fiber.Ctx, db *sql.DB) error {
	var cars []Car

	if err := c.BodyParser(&cars); err != nil {
		log.Warn(err)
		return c.Status(fiber.StatusBadRequest).JSON("Error parsing request body")
	}

	// Potential future improvement to accept valid cars and reject invalid ones
	for _, car := range cars {
		if err := carValidation(car); err != nil {
			log.Warn(err)
			errorString := fmt.Sprintf("%v", err)
			return c.Status(fiber.StatusBadRequest).JSON(errorString)
		}
	}

	for _, car := range cars {
		_, err := db.Exec("INSERT into cars (make, model, builddate, colourid) VALUES (?, ?, ?, ?)", car.Make, car.Model, car.BuildDate, car.ColourID)
		if err != nil {
			log.Warn(err)
			return c.Status(fiber.StatusInternalServerError).JSON("Error saving car to database")
		}
	}

	return c.SendStatus(fiber.StatusCreated)
}

func getCarByIdHandler(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	if id == "" {
		log.Warn("Missing param id on /car/:id")
		return c.Status(fiber.StatusBadRequest).JSON("Parameter 'id' cannot be empty")
	}

	var car Car
	row := db.QueryRow("SELECT * from cars WHERE id = ?", id)
	row.Scan(&car.ID, &car.Make, &car.Model, &car.BuildDate, &car.ColourID)
	return c.JSON(car)
}

func deleteCarByIdHandler(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	if id == "" {
		log.Warn("Missing param id on /car/:id")
		return c.Status(fiber.StatusBadRequest).JSON("Parameter 'id' cannot be empty")
	}

	_, err := db.Exec("DELETE from cars WHERE id = ?", id)
	if err != nil {
		log.Warn(err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error deleting vehicle")
	}
	return c.Status(fiber.StatusOK).JSON("Vehicle deleted successfully")
}

// Validation functions - can be combined later to take an interface{} as they're both very similar
func carValidation(c Car) error {
	validate = validator.New()
	err := validate.Struct(c)
	if err != nil {
		return err
	}
	return nil
}

func colourValidation(c Colour) error {
	validate = validator.New()
	err := validate.Struct(c)
	if err != nil {
		return err
	}
	return nil
}
