package utility

import (
	"letstalk/server/core/secrets"

	"github.com/namsral/flag"
)

// Flags to get db credentials
var (
	dbUser = flag.String("db_user", "", "mySQL user")
	dbPass = flag.String("db_pass", "", "mySQL password")
	dbAddr = flag.String("db_addr", "", "address of the database connection")

	secretsPath = flag.String("secrets_path", "~/secrets.json", "path to secrets.json")
)

// Flags to get es credentials
var (
	esAddr = flag.String("es_addr", "", "address of the elasticsearch connection")
)

// Methods to run before a client is initialized
func Bootstrap() {
	flag.Parse()

	// bootstrap secrets from local file
	secrets.LoadSecrets(*secretsPath)
}
