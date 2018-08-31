package main

import (
	"context"

	"letstalk/server/core/search"
	"letstalk/server/data"
	"letstalk/server/utility"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/namsral/flag"
)

// This job indexes multi traits for all user cohorts, positions and simple traits.
func main() {
	flag.Parse()

	db, err := utility.GetDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	es, err := utility.GetES()
	if err != nil {
		panic(err)
	}

	ids := make([]string, 0)
	multiTraits := make([]interface{}, 0)

	userSimpleTraits := make([]data.UserSimpleTrait, 0)
	err = db.Find(&userSimpleTraits).Error
	if err != nil {
		panic(err)
	}

	for _, trait := range userSimpleTraits {
		id, multiTrait := search.NewMultiTraitFromUserSimpleTrait(&trait)
		ids = append(ids, id)
		multiTraits = append(multiTraits, multiTrait)
	}

	userPositions := make([]data.UserPosition, 0)
	err = db.Find(&userPositions).Error
	if err != nil {
		panic(err)
	}

	for _, pos := range userPositions {
		id, multiTrait := search.NewMultiTraitFromUserPosition(&pos)
		ids = append(ids, id)
		multiTraits = append(multiTraits, multiTrait)
	}

	searchClient := search.NewClientWithContext(es, context.Background())
	err = searchClient.BulkIndexMultiTraits(ids, multiTraits)
	if err != nil {
		panic(err)
	}
}
