package main

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
)

// Quick and dirty utility function to reset the DB. Better version would be to read from a schema file.
func ResetDB() {
	db, err := ConnectToDb()
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	db.Exec(`DROP TABLE IF EXISTS "cars";
	CREATE TABLE IF NOT EXISTS "cars" (
		"id"	INTEGER NOT NULL UNIQUE,
		"make"	TEXT NOT NULL,
		"model"	TEXT NOT NULL,
		"builddate"	TEXT NOT NULL,
		"colourid"	INTEGER NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	DROP TABLE IF EXISTS "colours";
	INSERT INTO "cars" VALUES (1,'Mercedes','A Class','2022-04-04',1);
	CREATE TABLE IF NOT EXISTS "colours" (
		"id"	INTEGER NOT NULL UNIQUE,
		"name"	TEXT NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	INSERT INTO "colours" VALUES (1,'red');
	INSERT INTO "colours" VALUES (2,'blue');
	INSERT INTO "colours" VALUES (3,'white');
	INSERT INTO "colours" VALUES (4,'black');
	COMMIT;`)
	fmt.Println("Database reset successfully")
}

func ConnectToDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return db, err
	}

	return db, nil
}

func ClearDbData() {
	db, err := ConnectToDb()
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	db.Exec(`DELETE from cars;DELETE from colours;`)
	fmt.Println("Database emptied successfully")
}
