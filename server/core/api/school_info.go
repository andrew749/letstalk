package api

// School information
// (currently optional since we should be able to create a user without this information intially)
type SchoolInfo struct {
	Program        string `json:"program"`
	Sequence       string `json:"sequence"`
	GraduatingYear int    `json:"grad_year"`
}
