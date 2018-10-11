package api

import (
	"letstalk/server/data"
	"time"
)

type ConnectionRequest struct {
	UserId        data.TUserID    `json:"userId" binding:"required"`
	IntentType    data.IntentType `json:"intentType" binding:"required"`
	SearchedTrait *string         `json:"searchedTrait"`
	Message       *string         `json:"message"`
	// TODO(aklen): add field for connection request message
	// Output fields
	CreatedAt  time.Time  `json:"createdAt"`
	AcceptedAt *time.Time `json:"acceptedAt"`
}

type AcceptConnectionRequest struct {
	UserId data.TUserID `json:"userId" binding:"required"`
}

type RemoveConnection struct {
	UserId data.TUserID `json:"userId" binding:"required"`
}

type CreateMentorshipType string
const (
	CREATE_MENTORSHIP_TYPE_DRY_RUN CreateMentorshipType = "DRY_RUN"
	CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN CreateMentorshipType = "NOT_DRY_RUN"
)

type CreateMentorshipByEmail struct {
	MenteeEmail string               `json:"menteeEmail" binding:"required"`
	MentorEmail string               `json:"mentorEmail" binding:"required"`
	RequestType CreateMentorshipType `json:"requestType" binding:"required"`
}
