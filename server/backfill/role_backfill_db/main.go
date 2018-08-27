package main

import (
	"encoding/csv"
	"os"

	"letstalk/server/data"
	"letstalk/server/utility"

	"github.com/jinzhu/gorm"
	"github.com/namsral/flag"
)

var (
	inFile = flag.String("in", "", "Input csv file containing roles and header")
)

// This job backfills roles from a csv which contains rows with only the name of the role.
// If a role with a particular name already exists, it will not create a new one with the same name.
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

	roles := make([]data.Role, 0, len(records))
	for i, record := range records {
		if i == 0 {
			// Skip first record
			continue
		}
		roles = append(roles, data.Role{Name: record[0], IsUserGenerated: false})
	}

	// TODO: Look into bulk upserts, but for now this should be fine
	tx := db.Begin()
	for _, role := range roles {
		var existingRole data.Role
		err := tx.Where(&data.Role{Name: role.Name}).First(&existingRole).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			tx.Rollback()
			panic(err)
		} else if err == nil {
			role.Id = existingRole.Id
		}
		err = tx.Save(&role).Error
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
