package data

type TUserSimpleTraitID EntID

// Stores all of the simple traits for a user.
type UserSimpleTrait struct {
	Id                     TUserSimpleTraitID `gorm:"primary_key;not null;auto_increment:true"`
	UserId                 TUserID            `gorm:"not null"`
	SimpleTraitId          TSimpleTraitID     `gorm:"not null"`
	SimpleTraitName        string             `gorm:"not null"` // Denormalized
	SimpleTraitType        SimpleTraitType    `gorm:"not null"` // Denormalized
	SimpleTraitIsSensitive bool               `gorm:"not null"` // Denormalized
	Times
	// Untested
	User        *User        `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	SimpleTrait *SimpleTrait `gorm:"foreignkey:SimpleTraitId;association_foreignkey:Id"`
}

func (trait *UserSimpleTrait) GetUser() *User {
	return trait.User
}
