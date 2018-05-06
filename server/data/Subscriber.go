package data

type Subscriber struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	ClassYear   int    `json:"classYear" binding:"required"`
	ProgramName string `json:"programName" binding:"required"`
	Email       string `json:"emailAddress" binding:"required" gorm:"primary key"`
}
