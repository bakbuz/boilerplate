package database

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	//_ "github.com/lib/pq"
	//_ "modernc.org/sqlite"
	_ "github.com/microsoft/go-mssqldb"
)

type DB struct {
	*sql.DB
}

// New ...
// Bir veritabanı bağlantısı açar
func New(connectionString string) (*DB, error) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(3)

	return &DB{db}, nil
}

func GetDatabaseName(connectionString string) string {
	var dbName string
	var strs = strings.Split(connectionString, ";")
	for _, str := range strs {
		if strings.HasPrefix(strings.TrimSpace(str), "database") {
			dbName = strings.Split(str, "=")[1]
			break
		}
	}
	return dbName
}

func (db *DB) GetInt(query string) (int, error) {
	row := db.QueryRow(query)
	if row.Err() != nil {
		return -1, row.Err()
	}

	var value int
	if err := row.Scan(&value); err != nil {
		if err == sql.ErrNoRows { // sql: no rows in result set
			return -1, errors.New("no rows")
		}
		return -1, err
	}
	return value, nil
}

func (db *DB) GetString(query string) (string, error) {
	row := db.QueryRow(query)
	if row.Err() != nil {
		return "", row.Err()
	}

	var str string
	if err := row.Scan(&str); err != nil {
		if err == sql.ErrNoRows { // sql: no rows in result set
			return "", errors.New("no rows")
		}
		return "", err
	}

	return str, nil
}

func (db *DB) Count(ctx context.Context, query string) (int64, error) {
	var count int64

	row := db.QueryRowContext(ctx, query)
	if err := row.Scan(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func (db *DB) Get(result interface{}, query string, args ...any) error {
	row := db.QueryRow(query, args)

	err := row.Scan(result)

	return err
}
