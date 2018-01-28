package data_layer

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DBHelper struct {
	ConnectionPool sql.DB
}

type IDBHelper interface {
}

const (
	READ_USER       = "normal"
	READ_PASSWORD   = "cirocinhint" // FIXME: remove (this isn't the password)
	DATABASE_SERVER = "127.0.0.1"
	DATABASE_PORT   = "3306"
)

// FIXME: do some sort of connection pooling
func GetConnection() *sql.DB {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)",
		READ_USER,
		READ_PASSWORD,
		DATABASE_SERVER,
		DATABASE_PORT,
	)
	db, err := sql.Open(
		"mysql",
		connectionString,
	)
	if err != nil {
		log.Fatal("Unable to open db connection")
	}
	return db
}
