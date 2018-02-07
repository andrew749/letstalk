package main

import (
	"log"
	"net/http"

	"flag"
	"fmt"
	"letstalk/server/core/routes"
	"letstalk/server/core/secrets"
	"letstalk/server/core/sessions"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mijia/modelq/gmq"
	"github.com/romana/rlog"
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
		rlog.Error(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		rlog.Error("failed to connect to database: ", err)
	}
	sessionManager := sessions.CreateSessionManager()

	router := routes.Register(db, &sessionManager)
	secrets.LoadSecrets(*secretsPath)
	// Start server
	rlog.Info("Serving on port 8080...")
	err = http.ListenAndServe(":8080", router)

	if err != nil {
		log.Fatal(err)
	}
}
