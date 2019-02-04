package data

import "database/sql/driver"

type TOrganizationID EntID

type OrganizationType string

func (u *OrganizationType) Scan(value interface{}) error { *u = OrganizationType(value.([]byte)); return nil }
func (u OrganizationType) Value() (driver.Value, error) { return string(u), nil }

const (
	ORGANIZATION_TYPE_COMPANY      OrganizationType = "COMPANY"
	ORGANIZATION_TYPE_CLUB         OrganizationType = "CLUB"
	ORGANIZATION_TYPE_SPORTS_TEAM  OrganizationType = "SPORTS_TEAM"
	ORGANIZATION_TYPE_UNDETERMINED OrganizationType = "UNDETERMINED"
)

var ALL_ORGANIZATION_TYPES = map[OrganizationType]interface{}{
	ORGANIZATION_TYPE_COMPANY:      nil,
	ORGANIZATION_TYPE_CLUB:         nil,
	ORGANIZATION_TYPE_SPORTS_TEAM:  nil,
	ORGANIZATION_TYPE_UNDETERMINED: nil,
}

// This table stores all available organizations, where an organization is a company, club or sports
// team that a user can be a part of. E.g. Facebook. Used mainly by the user_position (position
// traits) table to get information about available organizations (to show in UI) and to get the
// names of organizations by ID (since they need to be denormalized into the user_position table).
type Organization struct {
	Id              TOrganizationID  `gorm:"primary_key;not null;auto_increment:true"`
	Name            string           `gorm:"not null"`
	Type            OrganizationType `gorm:"not null"`
	IsUserGenerated bool             `gorm:"not null"`
}
