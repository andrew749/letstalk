package data

type TRoleID EntID

// This table stores all available roles, where a role is a position that a user can take on at a
// company/club/sports team/etc. E.g. Software Engineer. Used mainly by the user_position (position
// traits) table to get information about available roles (to show in UI) and to get the names of
// roles by ID (since they need to be denormalized into the user_position table).
type Role struct {
	Id   TRoleID `gorm:"primary_key;not null;auto_increment:true"`
	Name string  `gorm:"not null"`
}
