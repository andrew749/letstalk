package data

type TMatchRoundMatchID EntID

type MatchRoundMatch struct {
	Id TMatchRoundMatchID `gorm:"primary_key;not null;auto_increment:true"`

	// Id of associated match round
	MatchRoundId TMatchRoundID `gorm:"not null"`

	// User one in the matching, which is the mentee
	MenteeUser   *User   `gorm:"foreignkey:MenteeUserId"`
	MenteeUserId TUserID `gorm:"not null"`

	// User two in the matching, which is the mentor
	MentorUser   *User   `gorm:"foreignkey:MentorUserId"`
	MentorUserId TUserID `gorm:"not null"`

	// Score calculated between users, assumes score is a float
	// Mainly used for debugging
	Score float32 `gorm:"not null"`
}
