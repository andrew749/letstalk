package main

import (
	"context"

	"letstalk/server/core/search"
	"letstalk/server/data"
	"letstalk/server/utility"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/namsral/flag"
)

// This job indexes all of the simple organizations currently in the organizations table in the db.
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

	dataOrganizations := make([]data.Organization, 0)
	err = db.Find(&dataOrganizations).Error
	if err != nil {
		panic(err)
	}

	organizations := make([]search.Organization, len(dataOrganizations))
	for i, dataOrganization := range dataOrganizations {
		organizations[i] = search.NewOrganizationFromDataModel(dataOrganization)
	}

	searchClient := search.NewClientWithContext(es, context.Background())
	err = searchClient.BulkIndexOrganizations(organizations)
	if err != nil {
		panic(err)
	}
}
