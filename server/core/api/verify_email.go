package api

type SendAccountVerificationEmailRequest struct {
	// UW email address to send the verification request to.
	Email string `json:"email" binding:"required"`
}

type VerifyEmailRequest struct {
	Id string `json:"id" binding:"required"`
}
