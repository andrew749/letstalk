package api

type ForgotPasswordChangeRequest struct {
	ForgotPasswordRequestId string `json:"requestId" binding:"required"`
	NewPassword             string `json:"password" binding:"required"`
}

type GenerateForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}
