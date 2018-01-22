package data_layer

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DBHelper struct {
	ConnectionPool sql.DB
}

type IDBHelper interface {
}

// FIXME: do some sort of connection pooling
func GetConnection() {
	db, err := sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/hello")

	if err != nil {
		log.Fatal("Unable to open db connection")
	}

	return db
}
