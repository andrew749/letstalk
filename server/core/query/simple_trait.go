package query

import (
	"context"
	"fmt"
	"strings"

	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/search"
	"letstalk/server/data"

	"github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic"
	"github.com/romana/rlog"
)

func getSimpleTrait(db *gorm.DB, traitId data.TSimpleTraitID) (*data.SimpleTrait, errs.Error) {
	var trait data.SimpleTrait
	err := db.Where(&data.SimpleTrait{Id: traitId}).First(&trait).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errs.NewRequestError(fmt.Sprintf("Simple trait with id %d not found", traitId))
		}
		return nil, errs.NewDbError(err)
	}
	return &trait, nil
}

func indexSimpleTrait(es *elastic.Client, trait data.SimpleTrait) {
	if es != nil {
		searchClient := search.NewClientWithContext(es, context.Background())
		searchTrait := search.NewSimpleTraitFromDataModel(trait)
		err := searchClient.IndexSimpleTrait(searchTrait)
		if err != nil {
			raven.CaptureError(err, nil)
			rlog.Error(err)
		}
	} else {
		rlog.Warn(fmt.Sprintf("Not indexing simple trait %s since no es provided", trait.Name))
	}
}

// Returns a simple trait with the given name or creates a new one if one doesn't already exist.
// TODO: Maybe make this take `isSensitive` so that user can specify that when creating a new
// user generated simple trait.
func getOrCreateSimpleTrait(
	db *gorm.DB,
	es *elastic.Client,
	name string,
) (*data.SimpleTrait, errs.Error) {
	var trait data.SimpleTrait

	err := ctx.WithinTx(db, func(db *gorm.DB) error {
		err := db.Where(&data.SimpleTrait{Name: name}).First(&trait).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				trait = data.SimpleTrait{
					Name:            name,
					Type:            data.SIMPLE_TRAIT_TYPE_UNDETERMINED,
					IsSensitive:     false,
					IsUserGenerated: true,
				}
				// Add trait if it doesn't already exist.
				if err := db.Create(&trait).Error; err != nil {
					return err
				}

				go indexSimpleTrait(es, trait)
			} else {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, errs.NewDbError(err)
	}

	return &trait, nil
}

func addUserSimpleTrait(db *gorm.DB, userId data.TUserID, trait data.SimpleTrait) errs.Error {
	var userTrait data.UserSimpleTrait
	// TODO: Trying using `WithinTx`
	tx := db.Begin()
	err := tx.Where(
		&data.UserSimpleTrait{UserId: userId, SimpleTraitId: trait.Id},
	).First(&userTrait).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		tx.Rollback()
		return errs.NewDbError(err)
	} else if err == nil {
		tx.Rollback()
		return errs.NewRequestError(fmt.Sprintf("You already have the trait \"%s\"", trait.Name))
	}

	userTrait = data.UserSimpleTrait{
		UserId:                 userId,
		SimpleTraitId:          trait.Id,
		SimpleTraitName:        trait.Name,
		SimpleTraitType:        trait.Type,
		SimpleTraitIsSensitive: trait.IsSensitive,
	}
	if dbErr := tx.Save(&userTrait).Error; dbErr != nil {
		tx.Rollback()
		return errs.NewDbError(dbErr)
	}

	if dbErr := tx.Commit().Error; dbErr != nil {
		return errs.NewDbError(dbErr)
	}
	return nil
}

func AddUserSimpleTraitById(
	db *gorm.DB,
	userId data.TUserID,
	traitId data.TSimpleTraitID,
) errs.Error {
	trait, err := getSimpleTrait(db, traitId)
	if err != nil {
		return err
	}
	return addUserSimpleTrait(db, userId, *trait)
}

func AddUserSimpleTraitByName(
	db *gorm.DB,
	es *elastic.Client,
	userId data.TUserID,
	name string,
) errs.Error {
	name = strings.TrimSpace(name)
	trait, err := getOrCreateSimpleTrait(db, es, name)
	if err != nil {
		return err
	}
	return addUserSimpleTrait(db, userId, *trait)
}

// TODO: Take userTraitId instead of traitId for consistency with userPosition
func RemoveUserSimpleTrait(
	db *gorm.DB,
	userId data.TUserID,
	traitId data.TSimpleTraitID,
) errs.Error {
	toDelete := data.UserSimpleTrait{UserId: userId, SimpleTraitId: traitId}
	err := db.Delete(&toDelete).Error
	if err != nil {
		return errs.NewDbError(err)
	}
	return nil
}
