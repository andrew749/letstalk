package main

import (
	"log"
	"net/http"

	"fmt"
	"letstalk/server/core/routes"
	"letstalk/server/core/secrets"
	"letstalk/server/core/sessions"
	"letstalk/server/data"

	"github.com/gin-contrib/pprof"
	"github.com/namsral/flag"

	"github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
	profiling   = flag.Bool("profiling", false, "Whether to turn on profiling endpoints.")
	production  = flag.Bool("PROD", false, "Whether to run in debug mode.")
)

func main() {
	rlog.Info("Starting server")
	flag.Parse()

	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf("%s:%s@%s/letstalk?charset=utf8mb4&parseTime=true", *dbUser, *dbPass, *dbAddr),
	)

	if err != nil {
		rlog.Error(err)
		panic("Failed to connect to database.")
	}

	defer db.Close()

	// log in development
	db.LogMode(!*production)

	// create the database
	data.CreateDB(db.Set("gorm:table_options", "CHARSET=utf8mb4")) // Create tables using utf8mb4 encoding. Only works with MySQL.

	sessionManager := sessions.CreateSessionManager(db)
	router := routes.Register(db, &sessionManager)
	if *profiling {
		// add cpu profiling
		pprof.Register(router, nil)
	}

	secrets.LoadSecrets(*secretsPath)

	// setup sentry logging
	raven.SetDSN(secrets.GetSecrets().SentryDSN)

	// production specific setup
	if *production {
		rlog.Info("Running in Production")
		raven.SetTagsContext(map[string]string{
			"environment": "production",
		})
		// setup sentry
	} else {
		rlog.Info("Running in Development mode")
		raven.SetTagsContext(map[string]string{
			"environment": "development",
		})
	}

	// Start server
	rlog.Info("Serving on port ", *port)

	// catch error and log
	raven.CapturePanic(func() {
		err = http.ListenAndServe(fmt.Sprintf(":%s", *port), router)
		if err != nil {
			raven.CaptureError(err, nil)
			log.Fatal(err)
		}
	}, nil)
}
