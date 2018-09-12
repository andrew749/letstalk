package utility

import (
	"letstalk/server/core/search"

	"github.com/olivere/elastic"
	"github.com/romana/rlog"
)

// Gets an elasticsearch client using command line params
func GetES() (*elastic.Client, error) {
	rlog.Infof("Connecting to elastic search instance at %s", *esAddr)
	return search.NewEsClient(*esAddr)
}
