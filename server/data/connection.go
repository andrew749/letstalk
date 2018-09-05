package data

import (
	"time"
)

type TConnectionID EntID

type Connection struct {
	ConnectionId TConnectionID `gorm:"not null;primary_key;auto_increment"`
	UserOne      *User         `gorm:"foreignkey:UserOneId"`
	UserOneId    TUserID
	UserTwo      *User         `gorm:"foreignkey:UserTwoId"`
	UserTwoId    TUserID
	CreatedAt    time.Time     `gorm:"not null"`
	UpdatedAt    time.Time
	DeletedAt    time.Time
	AcceptedAt   time.Time // Null until user two accepts.
}

type IntentType string
const (
	INTENT_TYPE_SEARCH IntentType = "search"
	INTENT_TYPE_REC_GENERAL IntentType = "recommendation_general"
	INTENT_TYPE_REC_COHORT IntentType = "recommendation_cohort"
)

type ConnectionIntent struct {
	ConnectionId  TConnectionID `gorm:"not null;primary_key"`
	Connection    *Connection   `gorm:"foreignkey:Connection"`
	Type          IntentType    `gorm:"not null;size:100"`
	SearchedTrait *string       `gorm:"type:text"` // Only applies to "search" type
}

type Mentorship struct {
	ConnectionId TConnectionID `gorm:"not null;primary_key"`
	Connection   *Connection   `gorm:"foreignkey:Connection"`
	MentorUser   *User         `gorm:"foreignkey:MentorUserId"`
	MentorUserId TUserID
	CreatedAt    time.Time     `gorm:"not null"`
	UpdatedAt    time.Time
	DeletedAt    time.Time
}
