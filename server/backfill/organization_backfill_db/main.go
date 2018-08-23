package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"letstalk/server/data"
	"letstalk/server/utility"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/namsral/flag"
)

var (
	inFile = flag.String("in", "", "Input csv file containing organizations and header")
)

// This job backfills organizations from a csv which contains rows with the following format
// "name,type". If an organization with a particular name already exists, it will not create a new
// one with the same name, but instead override it with new type.
// Note that the csv also contains a header, so the first row is skipped.
func main() {
	flag.Parse()

	db, err := utility.GetDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	f, err := os.Open(*inFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	organizations := make([]data.Organization, 0, len(records))
	for i, record := range records {
		if i == 0 {
			// Skip first record
			continue
		}
		tpe := data.OrganizationType(record[1])
		if _, ok := data.ALL_ORGANIZATION_TYPES[tpe]; !ok {
			panic(fmt.Sprintf("Record %d has an invalid type %s\n", i, record[1]))
		}

		organizations = append(organizations, data.Organization{
			Name:            record[0],
			Type:            tpe,
			IsUserGenerated: false,
		})
	}

	// TODO: Look into bulk upserts, but for now this should be fine
	tx := db.Begin()
	for _, organization := range organizations {
		var existingOrganization data.Organization
		err := tx.Where(&data.Organization{Name: organization.Name}).First(&existingOrganization).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			tx.Rollback()
			panic(err)
		} else if err == nil {
			organization.Id = existingOrganization.Id
		}
		err = tx.Save(&organization).Error
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	err = tx.Commit().Error
	if err != nil {
		panic(err)
	}
}
