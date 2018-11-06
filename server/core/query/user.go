package query

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"errors"

	"github.com/jinzhu/gorm"
	"letstalk/server/core/utility/uw_email"
)

// Returns missing users out of the given list
func MissingUsers(db *gorm.DB, userIds []data.TUserID) ([]data.TUserID, errs.Error) {
	users := make([]data.User, len(userIds))
	for i, userId := range userIds {
		users[i] = data.User{UserId: userId}
	}
	var foundUsers []data.User
	err := db.Where(&users).Find(&foundUsers).Error
	if err != nil {
		return nil, errs.NewDbError(err)
	}
	existingUserIds := make(map[data.TUserID]interface{})
	for _, user := range foundUsers {
		existingUserIds[user.UserId] = nil
	}
	missingUserIds := make([]data.TUserID, 0)
	for _, userId := range userIds {
		if _, ok := existingUserIds[userId]; !ok {
			missingUserIds = append(missingUserIds, userId)
		}
	}
	return missingUserIds, nil
}

func GetUserById(db *gorm.DB, userId data.TUserID) (*data.User, error) {
	var user data.User
	if db.Where(&data.User{UserId: userId}).First(&user).RecordNotFound() {
		return nil, errs.NewNotFoundError("Unable to find user")
	}
	return &user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*data.User, error) {
	var user data.User
	var err error
	result := db.Where(&data.User{Email: email}).First(&user)
	if result.Error != nil {
		return nil, errs.NewDbError(result.Error)
	}
	if !result.RecordNotFound() {
		return &user, nil
	}
	// Fallback to looking for uwaterloo email
	if !uw_email.Validate(email) {
		return nil, nil
	}
	var verifyEmailId *data.VerifyEmailId
	uwEmail := uw_email.FromString(email)
	if verifyEmailId, err = GetVerifyEmailIdByUwEmail(db, uwEmail); err != nil {
		return nil, errs.NewDbError(err)
	} else if verifyEmailId == nil {
		return nil, nil
	}
	if user, err := GetUserById(db, verifyEmailId.UserId); err != nil {
		return nil, errs.NewInternalError("Failed to get user: %v", err)
	} else {
		return user, nil
	}
}

func GetUserBySecret(db *gorm.DB, secret string) (*data.User, error) {
	var user data.User
	if db.Where(&data.User{Secret: secret}).First(&user).RecordNotFound() {
		return nil, errors.New("unable to find user")
	}
	return &user, nil
}

func GetUserProfileById(
	db *gorm.DB,
	userId data.TUserID,
	includeContactInfo bool,
) (*data.User, error) {
	var user data.User

	query := db.Where(
		&data.User{UserId: userId},
	).Preload("AdditionalData").Preload("UserPositions").Preload("UserSimpleTraits")

	if includeContactInfo {
		query = query.Preload("ExternalAuthData")
	}

	if query.First(&user).RecordNotFound() {
		return nil, errs.NewNotFoundError("Unable to find user")
	}

	return &user, nil
}

// Need to pass in all of this information just cause we want it to be a challenge to actaully
// delete a user.
// WARNING: DELETES EVERYTHING ABOUT A USER AND MAY HAVE PERMANENT EFFECTS.
func NukeUser(
	db *gorm.DB,
	email string,
	firstName string,
	lastName string,
	userId data.TUserID,
) error {
	return ctx.WithinTx(db, func(db *gorm.DB) error {
		var user data.User
		err := db.Where(&data.User{
			UserId:    userId,
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
		}).Find(&user).Error
		if err != nil {
			return err
		}

		// BEGIN TRAITS
		err = db.Where(&data.UserLocation{UserId: userId}).Delete(&data.UserLocation{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.UserSimpleTrait{UserId: userId}).Delete(&data.UserSimpleTrait{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.UserPosition{UserId: userId}).Delete(&data.UserPosition{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.UserCohort{UserId: userId}).Delete(&data.UserCohort{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.UserGroup{UserId: userId}).Delete(&data.UserGroup{}).Error
		if err != nil {
			return err
		}
		// END TRAITS

		// BEGIN EXTRA USER DATA
		err = db.Where(&data.AuthenticationData{UserId: userId}).Delete(&data.AuthenticationData{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.UserDevice{UserId: userId}).Delete(&data.UserDevice{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.UserAdditionalData{UserId: userId}).Delete(&data.UserAdditionalData{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.ExternalAuthData{UserId: userId}).Delete(&data.ExternalAuthData{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.Session{UserId: userId}).Delete(&data.Session{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.Notification{UserId: userId}).Delete(&data.Notification{}).Error
		if err != nil {
			return err
		}
		// END EXTRA USER DATA

		// BEGIN TOKENS
		err = db.Where(&data.FbAuthToken{UserId: userId}).Delete(&data.FbAuthToken{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.VerifyEmailId{UserId: userId}).Delete(&data.VerifyEmailId{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.ForgotPasswordId{UserId: userId}).Delete(&data.ForgotPasswordId{}).Error
		if err != nil {
			return err
		}
		// END TOKENS

		// BEGIN OLD STUFF
		err = db.Where(&data.RequestMatching{Asker: userId}).Delete(&data.RequestMatching{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.RequestMatching{Answerer: userId}).Delete(&data.RequestMatching{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.Matching{Mentor: userId}).Delete(&data.Matching{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.Matching{Mentee: userId}).Delete(&data.Matching{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.UserVector{UserId: userId}).Delete(&data.UserVector{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.UserCredential{UserId: userId}).Delete(&data.UserCredential{}).Error
		if err != nil {
			return err
		}
		err = db.Where(&data.UserCredentialRequest{UserId: userId}).Delete(
			&data.UserCredentialRequest{},
		).Error
		if err != nil {
			return err
		}
		// END OLD STUFF

		connections, err := GetAllConnections(db, userId)
		if err != nil {
			return err
		}
		for _, connection := range connections {
			if connection.Intent != nil {
				err = db.Delete(connection.Intent).Error
				if err != nil {
					return err
				}
			}
			if connection.Mentorship != nil {
				err = db.Delete(connection.Mentorship).Error
				if err != nil {
					return err
				}
			}
			err = db.Delete(&connection).Error
			if err != nil {
				return err
			}
		}

		return db.Delete(&user).Error
	})
}
