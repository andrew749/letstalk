package main

import (
	"log"
	"net/http"

	"flag"
	"fmt"
	"letstalk/server/core/routes"
	"letstalk/server/core/secrets"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mijia/modelq/gmq"
)

// DB flags
var (
	dbUser = flag.String("db-user", "", "mySQL user")
	dbPass = flag.String("db-pass", "", "mySQL password")
	dbAddr = flag.String("db-addr", "", "address of the database connection")
)

// Authentication flags
var (
	secretsPath = flag.String("secrets-path", "~/secrets.json", "path to secrets.json")
)

func main() {
	flag.Parse()
	db, err := gmq.Open("mysql", fmt.Sprintf("%s:%s@%s/letstalk", *dbUser, *dbPass, *dbAddr))
	if err != nil {
		log.Print(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Print("failed to connect to database: ", err)
	}

	router := routes.Register(db)
	secrets.LoadSecrets(*secretsPath)
	// Start server
	http.ListenAndServe(":8080", router)
}
