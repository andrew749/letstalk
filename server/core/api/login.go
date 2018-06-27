package api

import "time"

type LoginRequestData struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	// optional token to associate with this session
	NotificationToken *string `json:"notificationToken"`
}

type LoginResponse struct {
	SessionId  string    `json:"sessionId"`
	ExpiryDate time.Time `json:"expiry"`
}
