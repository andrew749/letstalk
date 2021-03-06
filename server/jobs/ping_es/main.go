package main

import (
	"context"
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
	esClient, err := search.NewEsClient(*esEndpoint)
	if err != nil {
		rlog.Error(err)
		panic(err)
	}
	rlog.Debug("Connected to elastic search")
	// c := ctx.Context{Es: esClient}
	// context := context.Context{}
	client := search.NewClientWithContext(esClient, context.TODO())

	client.PrintAllSimpleTraits()
	client.CompletionSuggestionSimpleTraits("test", 10)
}
