package main

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
)

// Quick and dirty utility function to reset the DB. Better version would be to read from a schema file.
func ResetDB(dbPath string) {
	db, err := ConnectToDb(dbPath)
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
	fmt.Printf("Database %v reset successfully\n", dbPath)
}

func ConnectToDb(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return db, err
	}
	log.Infof("Connected to DB at %v", dbPath)
	return db, nil
}

func ClearDbData(dbPath string) {
	db, err := ConnectToDb(dbPath)
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	db.Exec(`DELETE from cars;DELETE from colours;`)
	fmt.Printf("Database %v emptied successfully\n", dbPath)
}
