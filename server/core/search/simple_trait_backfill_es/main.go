package main

import (
	"letstalk/server/core/search"
	"letstalk/server/data"
	"letstalk/server/utility"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/namsral/flag"
)

// This job indexes all of the simple traits currently in the simple_traits table in the db.
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

	dataTraits := make([]data.SimpleTrait, 0)
	err = db.Find(&dataTraits).Error
	if err != nil {
		panic(err)
	}

	traits := make([]search.SimpleTrait, len(dataTraits))
	for i, dataTrait := range dataTraits {
		traits[i] = search.NewSimpleTraitFromDataModel(dataTrait)
	}

	err = search.BulkIndexSimpleTraits(es, traits)
	if err != nil {
		panic(err)
	}
}
