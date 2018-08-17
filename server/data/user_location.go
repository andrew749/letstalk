package data

type TUserLocationID EntID

// Refers to the Google `place_id` for their places API
type TPlaceID EntID

type LocationType string

const (
	LOCATION_TYPE_HOMETOWN    LocationType = "HOMETOWN"
	LOCATION_TYPE_COOP_TERM   LocationType = "COOP_TERM"
	LOCATION_TYPE_SCHOOL_TERM LocationType = "SCHOOL_TERM"
)

// Stores Google places for users
type UserLocation struct {
	Id        TUserLocationID `gorm:"primary_key;not null;auto_increment:true"`
	User      User            `gorm:"foreignkey:UserId"`
	UserId    TUserID         `gorm:"not null"`
	PlaceId   TPlaceID        `gorm:"not null"`
	PlaceName string          `gorm:"not null"` // Plaintext description of the place (e.g. Waterloo, Ontario, Canada)
	Type      LocationType    `gorm:"not null"`
	StartDate string          `gorm:"not null"` // YYYY-MM-DD
	EndDate   string          // YYYY-MM-DD (optional)
	Times
}
