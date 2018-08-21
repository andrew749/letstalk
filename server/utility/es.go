package utility

import (
	"letstalk/server/core/search"

	"github.com/olivere/elastic"
)

// Gets an elasticsearch client using command line params
func GetES() (*elastic.Client, error) {
	return search.NewEsClient(*esAddr)
}
