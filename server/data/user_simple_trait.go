package data

type TUserSimpleTraitID EntID

// Stores all of the simple traits for a user.
type UserSimpleTrait struct {
	Id              TUserSimpleTraitID `gorm:"primary_key;not null;auto_increment:true"`
	SimpleTraitId   TSimpleTraitID     `gorm:"not null"`
	SimpleTraitName string             `gorm:"not null"` // Denormalized
	SimpleTraitType SimpleTraitType    `gorm:"not null"` // Denormalized
	IsSensitive     bool               `gorm:"not null"` // Denormalized
	Times
	// Untested
	SimpleTrait *SimpleTrait `gorm:"foreignkey:SimpleTraitId;association_foreignkey:Id"`
}
