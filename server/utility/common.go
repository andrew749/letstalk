package utility

import (
	"fmt"
	"letstalk/server/core/secrets"

	"github.com/namsral/flag"
)

var (
	// Flags to get db credentials
	dbUser = flag.String("db_user", "", "mySQL user")
	dbPass = flag.String("db_pass", "", "mySQL password")
	dbAddr = flag.String("db_addr", "", "address of the database connection")

	secretsPath = flag.String("secrets_path", "~/secrets.json", "path to secrets.json")
	isProd      = flag.Bool("PROD", false, "Whether to run in debug mode.")

	// Flags to get es credentials
	esAddr = flag.String("es_addr", "", "address of the elasticsearch connection")
)

var (
	bootstrapRun = false
)

var (
	QueueHelper SQSQueue
)

// Methods to run before a client is initialized
func Bootstrap() {
	if bootstrapRun {
		return
	}
	bootstrapRun = true
	flag.Parse()

	// bootstrap secrets from local file
	secrets.LoadSecrets(*secretsPath)

	// Setup a worker to process all messages sent to a queue.
	QueueHelper = CreateLocalSQSClient()
}

func checkBootstrapped() {
	if !bootstrapRun {
		Bootstrap()
	}
}

func IsProductionEnvironment() bool {
	checkBootstrapped()
	return *isProd
}

func GetWebsiteUrl() string {
	if IsProductionEnvironment() {
		return "https://hiveapp.org"
	}
	return "http://localhost"
}

func GetWebappUrl() string {
	return fmt.Sprintf("%s/web", GetWebsiteUrl())
}

func GetDeeplinkPrefix() string {
	checkBootstrapped()
	return secrets.GetSecrets().DeeplinkPrefix
}
