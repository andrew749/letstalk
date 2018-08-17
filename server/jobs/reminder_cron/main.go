package main

import (
	"letstalk/server/data"
	"letstalk/server/utility"

	"github.com/jinzhu/gorm"
)

func main() {
	utility.RunWithDb(func(tx *gorm.DB) error {
		var res []data.Matching
		if err := tx.Where("").Find(&res); err != nil {

		}
		return nil
	})
}
