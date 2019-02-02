package api

import "letstalk/server/data"

type VerifyLinkRequest struct {
	VerifyLinkId data.TVerifyLinkID `json:"verifyLinkId" binding:"required"`
}
