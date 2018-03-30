package data

// TODO: Fill function for this and run migration for it on start
type Program struct {
	ProgramId string `json:"programId" gorm:"not null;primary_key"`
}
