package api

type SubscriptionRequest struct {
	ClassYear    int    `json:"classYear" binding:"required"`
	ProgramName  string `json:"programName" binding:"required"`
	EmailAddress string `json:"emailAddress" binding:"required"`
	FirstName    string `json:"firstName" binding:"required"`
	LastName     string `json:"lastName" binding:"required"`
}

type SubscriptionResponse struct {
	Status string `json:"status"`
}
