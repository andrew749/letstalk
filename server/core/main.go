package main

import (
	"log"
	"net/http"

	"flag"
	"fmt"
	"letstalk/server/core/routes"
	"letstalk/server/core/secrets"

	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mijia/modelq/gmq"
)

// DB flags
var (
	dbUser = flag.String("db-user", "", "mySQL user")
	dbPass = flag.String("db-pass", "", "mySQL password")
	dbAddr = flag.String("db-addr", "", "address of the database connection")
)

// Auth flags
var (
	secretsPath = flag.String("secrets-path", "~/secrets.json", "path to secrets.json")
)

func main() {
	flag.Parse()
	db, err := gmq.Open("mysql", fmt.Sprintf("%s:%s@%s/letstalk", *dbUser, *dbPass, *dbAddr))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to database: ", err)
		os.Exit(1)
	}

	router := routes.Register(db)
	secrets.LoadSecrets(*secretsPath)
	// Start server
	http.ListenAndServe(":8080", router)
}
