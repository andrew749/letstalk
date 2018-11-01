package survey

import (
	"letstalk/server/core/api"
	"letstalk/server/data"
)

var allSurveysByGroup = map[data.SurveyGroup]api.Survey{
	Generic_v1.Group: Generic_v1,
}

func getSurveyDefinitionByGroup(group data.SurveyGroup) *api.Survey {
	if survey, ok := allSurveysByGroup[group]; ok {
		return &survey
	} else {
		return nil
	}
}
