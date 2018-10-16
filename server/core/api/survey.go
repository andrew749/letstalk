package api

import "letstalk/server/data"

type SurveyOption struct {
	Key  data.SurveyOptionKey `json:"key" binding:"required"`
	Text string          `json:"text" binding:"required"`
}

type SurveyQuestion struct {
	Key     data.SurveyQuestionKey `json:"key" binding:"required"`
	Prompt  string            `json:"prompt" binding:"required"`
	Options []SurveyOption    `json:"options" binding:"required"`
}

type Survey struct {
	Version   data.SurveyVersion   `json:"version" binding:"required"`
	Questions []SurveyQuestion     `json:"questions" binding:"required"`
	Responses data.SurveyResponses `json:"responses" binding:"required"`
}
