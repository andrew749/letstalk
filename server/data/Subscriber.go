package data

type Subscriber struct {
	FirstName   string `gorm:"size:190" json:"firstName" binding:"required"`
	LastName    string `gorm:"size:190" json:"lastName" binding:"required"`
	ClassYear   int    `gorm:"size:190" json:"classYear" binding:"required"`
	ProgramName string `gorm:"size:190" json:"programName" binding:"required"`
	Email       string `gorm:"size:190" json:"emailAddress" binding:"required" gorm:"primary key"`
}
