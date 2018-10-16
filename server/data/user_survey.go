package data

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type SurveyVersion string
type SurveyQuestionKey string
type SurveyOptionKey string

type SurveyResponses map[SurveyQuestionKey]SurveyOptionKey

type UserSurvey struct {
	gorm.Model
	UserId    TUserID         `gorm:"not null"`
	Version   SurveyVersion   `gorm:"not null;size:190"`
	Responses SurveyResponses `gorm:"not null;type:text"`
}

func (u *SurveyVersion) Scan(value interface{}) error { *u = SurveyVersion(value.([]uint8)); return nil }
func (u SurveyVersion) Value() (driver.Value, error)  { return string(u), nil }
func (u *SurveyQuestionKey) Scan(value interface{}) error { *u = SurveyQuestionKey(value.([]uint8)); return nil }
func (u SurveyQuestionKey) Value() (driver.Value, error)  { return string(u), nil }
func (u *SurveyOptionKey) Scan(value interface{}) error { *u = SurveyOptionKey(value.([]uint8)); return nil }
func (u SurveyOptionKey) Value() (driver.Value, error)  { return string(u), nil }

func (s *SurveyResponses) Scan(value interface{}) error {
	reader := bytes.NewReader(value.([]byte))
	return json.NewDecoder(reader).Decode(s)
}

func (s SurveyResponses) Value() (driver.Value, error)  {
	buf := bytes.Buffer{}
	json.NewEncoder(&buf).Encode(s)
	return string(buf.Bytes()), nil
}
