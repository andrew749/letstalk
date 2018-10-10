package api

type FBLoginRequestDataCore struct {
	Token  string `json:"token" binding:"required"`
	Expiry int64  `json:"expiry" binding:"required"`
}

/**
 * Login with fb
 */
type FBLoginRequestData struct {
	FBLoginRequestDataCore
	NotificationToken *string `json:"notificationToken"`
}
