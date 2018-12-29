package main

import (
	"context"
	"log"
	"net/http"

	"fmt"
	"letstalk/server/constants"
	"letstalk/server/core/errs"
	"letstalk/server/core/routes"
	"letstalk/server/core/search"
	"letstalk/server/core/secrets"
	"letstalk/server/core/sessions"
	"letstalk/server/data"
	sqs_notification_processor "letstalk/server/jobs/sqs_notification_processor/src"

	"github.com/aws/aws-lambda-go/events"
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

	var es *elastic.Client

	if *useElastic {
		es, err = utility.GetES()
		if err != nil {
			rlog.Error(err)
			panic("Failed to connect to elasticsearch.")
		}
		rlog.Infof("Elastic search status: %s", es.String())
		rlog.Info("Creating indexes in ES")
		searchClient := search.NewClientWithContext(es, context.Background())
		if err := searchClient.CreateEsIndexes(); err != nil {
			// Failures here are okay since the indexes could already exist.
			rlog.Errorf("Error creating indexes: %#v", err)
		} else {
			rlog.Info("Success creating indexes in ES")
		}
	}

	// log in development
	db.LogMode(!utility.IsProductionEnvironment())

	// Create tables using utf8mb4 encoding. Only works with MySQL.
	data.CreateDB(db.Set("gorm:table_options", "CHARSET=utf8mb4"))

	sessionManager := sessions.CreateSessionManager(db)
	router := routes.Register(db, es, &sessionManager)
	if *profiling {
		// add cpu profiling
		pprof.Register(router, nil)
	}

	// setup sentry logging
	raven.SetDSN(secrets.GetSecrets().SentryDSN)

	// production specific setup
	if utility.IsProductionEnvironment() {
		rlog.Info("Running in Production mode")
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

	// setup queue listeners for local delivery
	sqs := utility.QueueHelper.(utility.LocalQueueImpl)
	go sqs.QueueProcessor()
	sqs.SubscribeListener(constants.NotificationQueueUrl, func(event *events.SQSEvent) error {
		if err := sqs_notification_processor.SendNotificationLambda(*event); err != nil {
			rlog.Critical(err)
			return err
		}
		return nil
	})

	// finish processing any messages before closing
	defer sqs.WaitForQueueDone()
	defer sqs.CloseQueue()

	// Start server
	rlog.Info("Serving on port ", *port)

	// catch error and log
	raven.CapturePanic(func() {
		err = http.ListenAndServe(fmt.Sprintf(":%s", *port), router)
		if err != nil {
			raven.CaptureError(err, nil)
			log.Fatal(err.(*errs.BaseError).VerboseError())
		}
	}, nil)
}
