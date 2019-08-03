package data

import (
	"database/sql/driver"
	"time"
)

type TConnectionID EntID

type Connection struct {
	ConnectionId TConnectionID `gorm:"not null;primary_key;auto_increment"`
	UserOne      *User         `gorm:"foreignkey:UserOneId"`
	UserOneId    TUserID
	UserTwo      *User `gorm:"foreignkey:UserTwoId"`
	UserTwoId    TUserID
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	AcceptedAt   *time.Time             // Null until user two accepts.
	Intent       *ConnectionIntent      `gorm:"foreignkey:ConnectionId"`
	Mentorship   *Mentorship            `gorm:"foreignkey:ConnectionId"`
	MatchRounds  []ConnectionMatchRound `gorm:"foreignkey:ConnectionId"`
}

type IntentType string

const (
	INTENT_TYPE_SEARCH      IntentType = "SEARCH"
	INTENT_TYPE_REC_GENERAL IntentType = "RECOMMENDATION_GENERAL"
	INTENT_TYPE_REC_COHORT  IntentType = "RECOMMENDATION_COHORT"
	INTENT_TYPE_SCAN_CODE   IntentType = "SCAN_CODE"
	INTENT_TYPE_ASSIGNED    IntentType = "ASSIGNED_MATCH"
)

func (u *IntentType) Scan(value interface{}) error { *u = IntentType(value.([]byte)); return nil }
func (u IntentType) Value() (driver.Value, error)  { return string(u), nil }

type ConnectionIntent struct {
	ConnectionId  TConnectionID `gorm:"not null;primary_key"`
	Type          IntentType    `gorm:"not null;size:100"`
	SearchedTrait *string       `gorm:"type:text"` // Only applies to "search" type
	Message       *string       `gorm:"type:text"`
}

type Mentorship struct {
	ConnectionId TConnectionID `gorm:"not null;primary_key"`
	MentorUser   *User         `gorm:"foreignkey:MentorUserId"`
	MentorUserId TUserID
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type ConnectionMatchRoundID EntID

type ConnectionMatchRound struct {
	Id           ConnectionMatchRoundID `gorm:"not null;primary_key"`
	ConnectionId TConnectionID          `gorm:"not null"`
	MatchRoundId TMatchRoundID          `gorm:"not null"`
}
