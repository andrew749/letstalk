package data

type Subscriber struct {
	FirstName   string `gorm:"size:100" json:"firstName" binding:"required"`
	LastName    string `gorm:"size:100" json:"lastName" binding:"required"`
	ClassYear   int    `gorm:"size:100" json:"classYear" binding:"required"`
	ProgramName string `gorm:"size:100" json:"programName" binding:"required"`
	Email       string `gorm:"size:100" json:"emailAddress" binding:"required" gorm:"primary key"`
}
