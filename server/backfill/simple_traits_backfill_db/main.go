package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"letstalk/server/data"
	"letstalk/server/utility"

	"github.com/jinzhu/gorm"
	"github.com/namsral/flag"
)

var (
	inFile = flag.String("in", "", "Input csv file containing simple traits and header")
)

// This job backfills simple traits from a csv with the row format "name,type,isSensitive". Will
// override simple traits with the same name, potentially updating their type and isSensitive.
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

	traits := make([]data.SimpleTrait, 0, len(records))
	for i, record := range records {
		if i == 0 {
			// Skip first record
			continue
		}
		var (
			isSensitive = true
			tpe         = data.SimpleTraitType(record[1])
		)
		if strings.ToUpper(strings.TrimSpace(record[2])) == "FALSE" {
			isSensitive = false
		}

		if _, ok := data.ALL_SIMPLE_TRAIT_TYPES[tpe]; !ok {
			panic(fmt.Sprintf("Record %d has an invalid type %s\n", i, record[1]))
		}

		traits = append(traits, data.SimpleTrait{
			Name:            record[0],
			Type:            tpe,
			IsSensitive:     isSensitive,
			IsUserGenerated: false,
		})
	}

	// TODO: Look into bulk upserts, but for now this should be fine
	tx := db.Begin()
	for _, trait := range traits {
		var existingTrait data.SimpleTrait
		err := tx.Where(&data.SimpleTrait{Name: trait.Name}).First(&existingTrait).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			tx.Rollback()
			panic(err)
		} else if err == nil {
			trait.Id = existingTrait.Id
		}
		err = tx.Save(&trait).Error
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
