package survey

import (
	"letstalk/server/core/api"
	"letstalk/server/data"
)

var allSurveysByGroup = map[data.SurveyGroup]api.Survey{
	Generic_v1.Group: Generic_v1,
}

var allSurveysByGroupId = map[data.TGroupID]api.Survey{
	"WICS": Generic_v1,
}

func GetSurveyDefinitionByGroup(group data.SurveyGroup) *api.Survey {
	if survey, ok := allSurveysByGroup[group]; ok {
		return &survey
	} else {
		return nil
	}
}

func GetSurveyDefinitionByGroupId(groupId data.TGroupID) *api.Survey {
	if survey, ok := allSurveysByGroupId[groupId]; ok {
		return &survey
	} else {
		return nil
	}
}
