package data

type TUserPositionID EntID

// One of the trait entity types. Describes positions that the user has held or currently holds.
type UserPosition struct {
	Id               TUserPositionID  `gorm:"primary_key;not null;auto_increment:true"`
	OrganizationId   TOrganizationID  `gorm:"not null"`
	OrganizationName string           `gorm:"not null"` // Denormalized
	OrganizationType OrganizationType `gorm:"not null"` // Denormalized
	RoleId           TRoleID          `gorm:"not null"`
	RoleName         string           `gorm:"not null"` // Denormalized
	StartDate        string           `gorm:"not null"` // YYYY-MM-DD
	EndDate          *string          // YYYY-MM-DD (optional)
	Times
	// Untested
	Organization *Organization `gorm:"foreignkey:OrganizationId;association_foreignkey:Id"`
	Role         *Role         `gorm:"foreignkey:RoleId;association_foreignkey:Id"`
}
