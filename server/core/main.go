package main

import (
	"net/http"
	"log"

	"github.com/mijia/modelq/gmq"
	_ "github.com/go-sql-driver/mysql"
	"letstalk/server/core/routes"
	"letstalk/server/core/secrets"
	"flag"
	"fmt"
)

// DB flags
var (
	dbUser = flag.String("db-user", "", "mySQL user")
	dbPass = flag.String("db-pass", "", "mySQL password")
	dbAddr = flag.String("db-addr", "", "address of the database connection")
)

func main() {
	flag.Parse()
	db, err := gmq.Open("mysql", fmt.Sprintf("%s:%s@%s/letstalk", *dbUser, *dbPass, *dbAddr))
	if err != nil {
		log.Fatal(err)
		return
	}
	router := routes.Register(db)
	secrets.GetSecrets()
	// Start server
	http.ListenAndServe(":8080", router)
}
