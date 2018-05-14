package api

type MeetingConfirmation struct {
	Secret string `json:"secret" binding:"required"`
}
