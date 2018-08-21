package utility

import (
	"github.com/namsral/flag"
)

// Flags to get db credentials
var (
	dbUser = flag.String("db_user", "", "mySQL user")
	dbPass = flag.String("db_pass", "", "mySQL password")
	dbAddr = flag.String("db_addr", "", "address of the database connection")
)

// Flags to get es credentials
var (
	esAddr = flag.String("es_addr", "", "address of the elasticsearch connection")
)