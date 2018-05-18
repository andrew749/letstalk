package data

import (
	"github.com/jinzhu/gorm"
)

type Credential struct {
	gorm.Model
	Name string `gorm:"not null"`
}
