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

	multiTraits := make(map[string]interface{})

	var (
		userCohorts      []data.UserCohort
		userSimpleTraits []data.UserSimpleTrait
		userPositions    []data.UserPosition
	)

	err = db.Preload("Cohort").Find(&userCohorts).Error
	if err != nil {
		panic(err)
	}

	for _, cohort := range userCohorts {
		id, multiTrait := search.NewMultiTraitFromUserCohort(&cohort)
		multiTraits[id] = multiTrait
	}

	err = db.Find(&userSimpleTraits).Error
	if err != nil {
		panic(err)
	}

	for _, trait := range userSimpleTraits {
		id, multiTrait := search.NewMultiTraitFromUserSimpleTrait(&trait)
		multiTraits[id] = multiTrait
	}

	err = db.Find(&userPositions).Error
	if err != nil {
		panic(err)
	}

	for _, pos := range userPositions {
		id, multiTrait := search.NewMultiTraitFromUserPosition(&pos)
		multiTraits[id] = multiTrait
	}

	searchClient := search.NewClientWithContext(es, context.Background())
	err = searchClient.BulkIndexMultiTraits(multiTraits)
	if err != nil {
		panic(err)
	}
}
