package data

type TCohortID EntID

type Cohort struct {
	CohortId     TCohortID `gorm:"not null;auto_increment;primary_key"`
	ProgramId    string    `gorm:"not null"`
	ProgramName  string    `gorm:"not null"`
	GradYear     uint      `gorm:"not null"`
	IsCoop       bool      `gorm:"not null"`
	SequenceId   *string
	SequenceName *string
}
