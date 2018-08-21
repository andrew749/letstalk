package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"letstalk/server/data"
	"letstalk/server/utility"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/namsral/flag"
)

var (
	inFile = flag.String("in", "", "Input csv file containing simple traits and header")
)

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

	// TODO: Look into bulk inserts, but for now this should be fine
	tx := db.Begin()
	for _, trait := range traits {
		err := tx.Create(&trait).Error
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
