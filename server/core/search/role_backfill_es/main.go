package main

import (
	"context"

	"letstalk/server/core/search"
	"letstalk/server/data"
	"letstalk/server/utility"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/namsral/flag"
)

// This job indexes all of the simple roles currently in the roles table in the db.
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

	dataRoles := make([]data.Role, 0)
	err = db.Find(&dataRoles).Error
	if err != nil {
		panic(err)
	}

	roles := make([]search.Role, len(dataRoles))
	for i, dataRole := range dataRoles {
		roles[i] = search.NewRoleFromDataModel(dataRole)
	}

	searchClient := search.NewClientWithContext(es, context.Background())
	err = searchClient.BulkIndexRoles(roles)
	if err != nil {
		panic(err)
	}
}
