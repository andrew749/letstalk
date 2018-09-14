package main

import (
	"letstalk/server/core/search"

	"github.com/namsral/flag"
	"github.com/romana/rlog"
)

var (
	esEndpoint = flag.String("es", "", "Elastic search endpoint to connect to")
)

func main() {
	flag.Parse()

	rlog.Debug("Connecting to elastic search")
	client, err := search.NewEsClient(*esEndpoint)
	if err != nil {
		rlog.Error(err)
		panic(err)
	}
	rlog.Debug("Connected to elastic search")
	rlog.Debug(client.ElasticsearchVersion(*esEndpoint))

	rlog.Info(client.String())
}
