package data

type TCohortID EntID

type Cohort struct {
	CohortId     TCohortID `gorm:"not null;auto_increment;primary_key"`
	ProgramId    string    `gorm:"not null;size:190"`
	ProgramName  string    `gorm:"not null;size:190"`
	GradYear     uint      `gorm:"not null"`
	IsCoop       bool      `gorm:"not null"`
	SequenceId   *string   `gorm:"size:190"`
	SequenceName *string   `gorm:"size:190"`
}
