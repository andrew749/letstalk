package main

import (
	"context"
	"log"
	"net/http"

	"fmt"
	"letstalk/server/core/routes"
	"letstalk/server/core/search"
	"letstalk/server/core/secrets"
	"letstalk/server/core/sessions"
	"letstalk/server/data"

	"github.com/gin-contrib/pprof"
	"github.com/namsral/flag"

	"letstalk/server/utility"

	"github.com/getsentry/raven-go"
	"github.com/olivere/elastic"
	"github.com/romana/rlog"
)

var (
	port = flag.String("port", "", "Port to host server on")
)

// Authentication flags
var (
	profiling  = flag.Bool("profiling", false, "Whether to turn on profiling endpoints.")
	isProd     = flag.Bool("PROD", false, "Whether to run in debug mode.")
	useElastic = flag.Bool("use_elastic", true, "Whether to create an Elasticsearch client")
)

func main() {
	rlog.Info("Starting server")
	utility.Bootstrap()

	db, err := utility.GetDB()

	if err != nil {
		rlog.Error(err)
		panic("Failed to connect to database.")
	}

	defer db.Close()

	var es *elastic.Client = nil

	// Right now, we never load the elasticsearch client on prod. This needs a little bit of infra
	// work.
	if *useElastic && !*isProd {
		es, err = utility.GetES()
		if err != nil {
			rlog.Error(err)
			panic("Failed to connect to elasticsearch.")
		}

		rlog.Info("Creating indexes in ES")
		searchClient := search.NewClientWithContext(es, context.Background())
		if err := searchClient.CreateEsIndexes(); err != nil {
			// Failures here are okay since the indexes could already exist.
			rlog.Error(err)
		} else {
			rlog.Info("Success creating indexes in ES")
		}
	}

	// log in development
	db.LogMode(!*isProd)

	// Create tables using utf8mb4 encoding. Only works with MySQL.
	if err := data.CreateDB(db.Set("gorm:table_options", "CHARSET=utf8mb4")); err != nil {
		panic(err)
	}

	sessionManager := sessions.CreateSessionManager(db)
	router := routes.Register(db, es, &sessionManager)
	if *profiling {
		// add cpu profiling
		pprof.Register(router, nil)
	}

	secrets.LoadSecrets(*secretsPath)

	// setup sentry logging
	raven.SetDSN(secrets.GetSecrets().SentryDSN)

	// production specific setup
	if *isProd {
		rlog.Info("Running in isProd")
		raven.SetTagsContext(map[string]string{
			"environment": "isProd",
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
