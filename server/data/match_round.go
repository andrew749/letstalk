package data

import (
	"database/sql/driver"
	"encoding/json"
)

type TMatchRoundID EntID

// Map of parameters used for the match round
type MatchParameters map[string]interface{}

// Unmarshall map from json
func (u *MatchParameters) Scan(value interface{}) error {
	var tmp map[string]interface{}
	err := json.Unmarshal(value.([]byte), &tmp)
	*u = tmp
	return err
}

// Marshall into json
func (u MatchParameters) Value() (driver.Value, error) {
	data, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return string(data), err
}

// Represents a round of matching for a particular group for mentor/mentee.
// A match round goes through the following states:
//                          ( match round created ) - matching algorithm generates matches
//                         /                       \
//            ( match round committed )       ( match round deleted )
//            - admin reviews and approves    - admin reviews and rejects matches
//            matches
//            - create CommitMatchRound job
//                        |
//               ( match round done)
//            - connections created and email sent
//            - admins of group notified of success
type MatchRound struct {
	Id TMatchRoundID `gorm:"primary_key;not null;auto_increment:true"`

	// Human-understandable name used to reference this match round for admins
	Name string `gorm:"not null"`

	// TODO(match-api): Figure out if this is the right ID after Andrew's changes
	GroupId TGroupID `gorm:"not null"`

	// Data mainly for debugging purposes since the parameters of the run aren't passed to jobmine
	MatchParameters MatchParameters `gorm:"type:text;not null"`

	// RunId for the CommitMatchRound job associated with this match round. If this is null, it means
	// that the job has not yet been created (i.e. admin has not yet committed the match round).
	RunId     *string    `gorm:"null"`
	CommitJob *JobRecord `gorm:"foreignkey:RunId;association_foreignkey:RunId"`

	// So that we can pull matches for this round
	Matches []MatchRoundMatch `gorm:"foreignkey:MatchRoundId;association_foreignkey:Id"`

	Times

	// TODO(match-api): Add association for the associated group
}
