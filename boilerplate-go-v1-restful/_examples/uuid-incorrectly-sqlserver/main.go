package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	db, err := sql.Open("sqlserver", "Server=127.0.0.1;Database=Catalog;User=sa;Password=123qwe..;TrustServerCertificate=True;MultipleActiveResultSets=true")
	if err != nil {
		log.Fatal("DEBUG 1: ", err)
	}
	defer db.Close()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(3)

	// Create Table
	var createCommand = `

IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'Catalog')
BEGIN
	CREATE DATABASE Catalog;
END;

IF NOT EXISTS (SELECT * FROM sys.sysobjects WHERE name = 'products' and xtype='U')
BEGIN
	CREATE TABLE products (
		[id] [uniqueidentifier] NOT NULL PRIMARY KEY,
		[name] [nvarchar](100) NOT NULL,
		[quantity] [int] NOT NULL,
		[created_at] [datetime2] NOT NULL
	)
END;
`

	_, err = db.Exec(createCommand)
	if err != nil {
		log.Fatal("DEBUG 2: ", err)
	}

	newid, _ := uuid.NewV7()
	var (
		id        uuid.UUID = newid
		name      string    = "Demo"
		quantity  int       = 1
		createdAt time.Time = time.Now().UTC()
	)

	// Insert
	_, err = db.Exec(`INSERT INTO products (id,name,quantity,created_at) VALUES (@Id,@Name,@Quantity,@CreatedAt)`,
		sql.Named("Id", id),
		sql.Named("Name", name),
		sql.Named("Quantity", quantity),
		sql.Named("CreatedAt", createdAt),
	)
	if err != nil {
		log.Fatal("DEBUG 3: ", err)
	}

	fmt.Println("INSERTED")
	fmt.Println(id, name, quantity, createdAt)

	// Get Last
	cursor, err := db.Query("SELECT TOP 1 * FROM products ORDER BY created_at DESC")
	if err != nil {
		log.Fatal("DEBUG 4: ", err)
	}
	defer cursor.Close()

	var (
		dbid        uuid.UUID
		dbname      string
		dbquantity  int
		dbcreatedAt time.Time
	)

	fmt.Println("FETCHED")

	for cursor.Next() {
		err = cursor.Scan(&dbid, &dbname, &dbquantity, &dbcreatedAt)
		if err != nil {
			log.Fatal("DEBUG 5: ", err)
		}

		fmt.Println(dbid, dbname, dbquantity, dbcreatedAt, id == dbid)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal("DEBUG 6: ", err)
	}
}
