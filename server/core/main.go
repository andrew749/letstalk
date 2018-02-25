package main

import (
	"log"
	"net/http"

	"fmt"
	"letstalk/server/core/routes"
	"letstalk/server/core/secrets"
	"letstalk/server/core/sessions"

	"github.com/namsral/flag"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mijia/modelq/gmq"
	"github.com/romana/rlog"
)

// DB flags
var (
	dbUser = flag.String("db_user", "", "mySQL user")
	dbPass = flag.String("db_pass", "", "mySQL password")
	dbAddr = flag.String("db_addr", "", "address of the database connection")
)

var (
	port = flag.String("port", "", "Port to host server on")
)

// Authentication flags
var (
	secretsPath = flag.String("secrets_path", "~/secrets.json", "path to secrets.json")
)

func main() {
	flag.Parse()

	db, err := gmq.Open("mysql", fmt.Sprintf("%s:%s@%s/letstalk", *dbUser, *dbPass, *dbAddr))
	if err != nil {
		rlog.Error(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		rlog.Error("failed to connect to database: ", err)
	}
	sessionManager := sessions.CreateSessionManager(db)

	router := routes.Register(db, &sessionManager)
	secrets.LoadSecrets(*secretsPath)
	// Start server
	rlog.Info("Serving on port ", *port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", *port), router)

	if err != nil {
		log.Fatal(err)
	}
}
