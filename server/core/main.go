package main

import (
	"log"
	"net/http"

	"fmt"
	"letstalk/server/core/routes"
	"letstalk/server/core/secrets"
	"letstalk/server/core/sessions"
	"letstalk/server/data"

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
)

func migrateDB(db *gorm.DB) {
	db.AutoMigrate(&data.AuthenticationData{})
	db.AutoMigrate(&data.Cohort{})
	data.PopulateCohort(db)
	db.AutoMigrate(&data.User{})
	db.AutoMigrate(&data.Session{})
	db.AutoMigrate(&data.UserVector{})
	db.AutoMigrate(&data.UserCohort{})
	db.AutoMigrate(&data.NotificationToken{})
	db.AutoMigrate(&data.FbAuthData{})
	db.AutoMigrate(&data.FbAuthToken{})
	db.AutoMigrate(&data.Matchings{})
}

func main() {
	rlog.Info("Starting server")
	flag.Parse()

	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf("%s:%s@%s/letstalk?parseTime=true", *dbUser, *dbPass, *dbAddr),
	)

	if err != nil {
		rlog.Error(err)
		panic("Failed to connect to database.")
	}

	defer db.Close()

	rlog.Info("Migrating database")
	db.LogMode(true)
	migrateDB(db)

	sessionManager := sessions.CreateSessionManager(db)

	router := routes.Register(db, &sessionManager)
	secrets.LoadSecrets(*secretsPath)

	// setup sentry
	raven.SetDSN(secrets.GetSecrets().SentryDSN)

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
