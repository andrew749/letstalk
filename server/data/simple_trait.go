package data

type TSimpleTraitID EntID

type SimpleTraitType string

const (
	SIMPLE_TRAIT_TYPE_INTEREST     SimpleTraitType = "INTEREST"
	SIMPLE_TRAIT_TYPE_EXPERIENCE   SimpleTraitType = "EXPERIENCE"
	SIMPLE_TRAIT_TYPE_RELIGION     SimpleTraitType = "RELIGION"
	SIMPLE_TRAIT_TYPE_RACE         SimpleTraitType = "RACE"
	SIMPLE_TRAIT_TYPE_UNDETERMINED SimpleTraitType = "UNDETERMINED"
)

// Stores all simple traits that can we written down in a single line of plaintext. Examples of
// these include experiences, hobbies, religion, race, etc. Since some of these are pretty
// sensitive, there is an isSensitive tag on each of these. These are mainly used by the
// user_simple_trait table to display a list of possible traits in the UI (when adding traits),
// and to pull the names of traits by ID for denormalization.
type SimpleTrait struct {
	Id              TSimpleTraitID  `gorm:"primary_key;not null;auto_increment:true"`
	Name            string          `gorm:"not null"`
	Type            SimpleTraitType `gorm:"not null"`
	IsSensitive     bool            `gorm:"not null"`
	IsUserGenerated bool            `gorm:"not null"`
}
