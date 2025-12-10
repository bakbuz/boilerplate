package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

func main() {

	// id, _ := uuid.NewV7()
	// fmt.Println(id)
	// return

	db, err := sql.Open("sqlite", "./testdb.sqlite")	
	if err != nil {
		log.Fatal("DEBUG 1: ", err)
	}
	defer db.Close()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(3)

	// Create Table
	var createCommand = `
CREATE TABLE IF NOT EXISTS "products" (
	"id"					TEXT NOT NULL UNIQUE,
	"name"				TEXT NOT NULL,
	"quantity"		INTEGER NOT NULL,
	PRIMARY KEY("id")
);`

	_, err = db.Exec(createCommand)
	if err != nil {
		log.Fatal("DEBUG 2: ", err)
	}

	var (
		id       uuid.UUID = uuid.MustParse("018e9ed9-9969-7dea-8dd0-2184bee56037") // uuid.NewV7()
		name     string    = "Demo"
		quantity int       = 1
	)

	// Insert
	_, err = db.Exec(`INSERT INTO products (id,name,quantity) VALUES (@Id,@Name,@Quantity)`,
		sql.Named("Id", id),
		sql.Named("Name", name),
		sql.Named("Quantity", quantity),
	)
	if err != nil {
		log.Fatal("DEBUG 3: ", err)
	}

	fmt.Println("INSERTED")
	fmt.Println(id, name, quantity)

	// GetAll
	cursor, err := db.Query("SELECT * FROM products")
	if err != nil {
		log.Fatal("DEBUG 4: ", err)
	}
	defer cursor.Close()

	var (
		dbid       uuid.UUID
		dbname     string
		dbquantity int
	)

	fmt.Println("FETCHED")

	for cursor.Next() {
		err = cursor.Scan(&dbid, &dbname, &dbquantity)
		if err != nil {
			log.Fatal("DEBUG 5: ", err)
		}

		fmt.Println(dbid, dbname, dbquantity, id == dbid)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal("DEBUG 6: ", err)
	}
}
