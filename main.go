package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	_ "github.com/mattn/go-sqlite3"
)

// Configuration
const prodDbPath = "file:./data/cardata.db?cache=shared"
const testDbPath = "file:./data/cardata_test.db?cache=shared"
const IP = "0.0.0.0"
const port = "8000"
const dateFormat = time.DateOnly
const buildDateMaxYears = 4

type Car struct {
	ID        int
	Make      string `validate:"required"`
	Model     string `validate:"required"`
	BuildDate string `validate:"required"`
	Colour    Colour `validate:"required"`
}

type CarAdd struct {
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
var clearDbFlag = flag.Bool("c", false, "clear database data on app start; overrides -r flag")
var resetDbFlag = flag.Bool("r", false, "reset database to default data on app start; is overridden by -c flag")
var useMetricsFlag = flag.Bool("m", false, "enable Fiber's metrics middleware")
var useCacheFlag = flag.Bool("cache", false, "enable Fiber's caching middleware")

func main() {
	flag.Parse()

	app := initApp(false) // splitting the app initialisation and listening out allows us to use app.Test for testing
	app.Listen(IP + ":" + port)
}

func initApp(testing bool) *fiber.App {
	var dbPath string
	if testing {
		log.Info("Testing mode enabled")
		dbPath = testDbPath
	} else {
		log.Info("Production mode enabled")
		dbPath = prodDbPath
	}

	if *clearDbFlag {
		ClearDbData(dbPath)
	} else if *resetDbFlag {
		ResetDB(dbPath)
	}

	db, err := ConnectToDb(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database at %v", dbPath)
	}

	app := fiber.New(fiber.Config{
		AppName:           "Drivvn API",
		EnablePrintRoutes: true,
		IdleTimeout:       2 * time.Minute,
	})

	if *useCacheFlag {
		log.Info("Enabled caching")
		app.Use(cache.New())
	} else {
		log.Info("Disabled caching")
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	if *useMetricsFlag {
		log.Info("Enabled metrics")
		app.Get("/metrics", monitor.New())
	} else {
		log.Info("Disabled metrics")
	}

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

func logRequest(c *fiber.Ctx) {
	log.Infof("%v - %v ", c.Method(), c.Path())
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	logRequest(c)
	return c.SendFile("readme.MD")
}

// GET /cars
func getCarsHandler(c *fiber.Ctx, db *sql.DB) error {
	logRequest(c)
	var cars []Car

	rows, err := db.Query("SELECT Car.id, Car.make, Car.model, Car.builddate, Colour.id, Colour.name FROM cars Car JOIN colours Colour ON Car.Colourid = Colour.id")
	if err != nil {
		log.Warn(err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error retrieving cars")
	}
	defer rows.Close()

	for rows.Next() {
		var car Car
		err := rows.Scan(&car.ID, &car.Make, &car.Model, &car.BuildDate, &car.Colour.ID, &car.Colour.Name)
		if err != nil {
			log.Warn(err)
			return c.Status(fiber.StatusInternalServerError).JSON("Error retrieving cars")
		}
		cars = append(cars, car)
	}

	if len(cars) == 0 {
		return c.Status(fiber.StatusNotFound).JSON("No cars found")
	}

	return c.JSON(cars)
}

// POST /cars
func postCarsHandler(c *fiber.Ctx, db *sql.DB) error {
	logRequest(c)
	var cars []CarAdd

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

		if err := colourValidation(car, db); err != nil {
			log.Warn(err)
			errorString := fmt.Sprintf("%v", err)
			return c.Status(fiber.StatusBadRequest).JSON(errorString)
		}

		if err := buildDateValidation(car.BuildDate); err != nil {
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

// GET /car/:id
func getCarByIdHandler(c *fiber.Ctx, db *sql.DB) error {
	logRequest(c)
	id := c.Params("id")
	if id == "" {
		log.Warn("Missing param id on /car/:id")
		return c.Status(fiber.StatusBadRequest).JSON("Parameter 'id' cannot be empty")
	}

	var car Car
	row := db.QueryRow("SELECT Car.id, Car.make, Car.model, Car.builddate, Colour.id, Colour.name FROM cars Car JOIN colours Colour ON Car.Colourid = Colour.id WHERE Car.id = ?", id)
	row.Scan(&car.ID, &car.Make, &car.Model, &car.BuildDate, &car.Colour.ID, &car.Colour.Name)

	if car.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON("Car not found with that id")
	}

	return c.JSON(car)
}

// DELETE /car/:id
func deleteCarByIdHandler(c *fiber.Ctx, db *sql.DB) error {
	logRequest(c)
	id := c.Params("id")
	if id == "" {
		log.Warn("Missing param id on /car/:id")
		return c.Status(fiber.StatusBadRequest).JSON("Parameter 'id' cannot be empty")
	}

	res, err := db.Exec("DELETE from cars WHERE id = ?", id)
	if err != nil {
		log.Warn(err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error deleting vehicle")
	}

	if recordsDeleted, _ := res.RowsAffected(); recordsDeleted < 1 {
		return c.Status(fiber.StatusNotFound).JSON("Vehicle not found with that id")
	}
	return c.Status(fiber.StatusOK).JSON("Vehicle deleted successfully")
}

// carValidation takes any interface and uses https://github.com/go-playground/validator to validate that fields are
// populated as required using the struct's tags. If validation fails it returns an error, otherwise it returns nil.
func carValidation(a interface{}) error {
	validate = validator.New()
	err := validate.Struct(a)
	if err != nil {
		return err
	}
	return nil
}

// colourValidation takes a Car struct and the current db connection
// It queries the DB for the given colour ID
// If no rows are returned by the query it returns an error, and otherwise returns nil
func colourValidation(c CarAdd, db *sql.DB) error {
	var colour Colour

	row := db.QueryRow("SELECT * FROM colours WHERE id = ?", c.ColourID)
	if err := row.Scan(&colour.ID, &colour.Name); err != nil {
		return errors.New("Colour validation failed - check ColourID exists")
	}

	return nil
}

// buildDateValidation takes a date string (in our chosen format of time.DateOnly), converts it to a time.Time
// It checks if it is more than the maximum allowed age (set by buildDateMaxYears)
// Returns an error if it is, and otherwise returns nil
func buildDateValidation(d string) error {
	formatString := dateFormat

	// convert string to date
	date, err := time.Parse(formatString, d)
	if err != nil {
		return err
	}

	buildAge := time.Since(date).Hours() / 24 / 365 // Not perfect to calculate calendar years

	if buildAge > buildDateMaxYears {
		errorString := fmt.Sprintf("Vehicle build date (%v) is older than the maximum allowed (%v years)", date, buildDateMaxYears)
		return errors.New(errorString)
	}
	return nil
}
