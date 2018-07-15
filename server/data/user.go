package data

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type UserRole string

const (
	USER_ROLE_DEFAULT UserRole = "DEFAULT"
	USER_ROLE_ADMIN   UserRole = "ADMIN"
)

// TUserID type of a user id
type TUserID EntID

type User struct {
	CreatedAt        time.Time `gorm:"not null"`
	UserId           TUserID   `gorm:"not null;primary_key;auto_increment"`
	FirstName        string    `gorm:"not null"`
	LastName         string    `gorm:"not null"`
	Email            string    `gorm:"type:varchar(128);not null;unique"`
	Secret           string    `gorm:"type:char(36);not null;unique"`
	Gender           int       `gorm:"not null"`
	Birthdate        string    `gorm:"type:varchar(100);not null"`
	Role             UserRole  `gorm:"not null"`
	ProfilePic       *string
	Sessions         []Session           `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	AuthData         *AuthenticationData `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	ExternalAuthData *ExternalAuthData   `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	Cohort           *UserCohort         `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	AdditionalData   *UserAdditionalData `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
}

func (u *UserRole) Scan(value interface{}) error { *u = UserRole(value.([]byte)); return nil }
func (u UserRole) Value() (driver.Value, error)  { return string(u), nil }

func CreateUser(
	db *gorm.DB,
	email string,
	firstName string,
	lastName string,
	gender int,
	birthdate string,
	role UserRole,
) (*User, error) {
	user := User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Gender:    gender,
		Birthdate: birthdate,
		Role:      role,
	}

	// Generate UUID for each user.
	secret, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	user.Secret = secret.String()

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
