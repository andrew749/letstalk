package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/data"
	"letstalk/server/utility"

	"github.com/namsral/flag"
)

var (
	inFile     = flag.String("in", "", "Input csv file containing emails of users to add")
	groupIdStr = flag.String("group_id", "", "Id for the group (uppercase please)")
	groupName  = flag.String("group_name", "", "Name for the group (capitalized please)")
)

func main() {
	flag.Parse()

	if groupIdStr == nil || *groupIdStr == "" {
		panic("Must provide -group_id")
	}
	if groupName == nil || *groupName == "" {
		panic("Must provide -group_name")
	}
	groupId := data.TGroupID(*groupIdStr)

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

	userIds := make([]data.TUserID, 0)
	missingEmails := make([]string, 0)
	for _, record := range records {
		user, err := query.GetUserByEmail(db, record[0])
		if err != nil {
			if _, ok := err.(*errs.NotFoundError); ok {
				missingEmails = append(missingEmails, record[0])
			} else {
				panic(err)
			}
		} else {
			userIds = append(userIds, user.UserId)
		}
	}

	if len(missingEmails) > 0 {
		fmt.Printf("Couldn't find the following users:\n")
		for _, email := range missingEmails {
			fmt.Printf("%s\n", email)
		}
		os.Exit(1)
	}

	if err := query.CreateUserGroups(db, userIds, groupId, *groupName); err != nil {
		panic(err)
	}
}
