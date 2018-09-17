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
