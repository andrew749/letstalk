package survey

import (
	"letstalk/server/core/api"
	"letstalk/server/data"
	"github.com/pkg/errors"
)

var allSurveysByGroup = map[data.SurveyGroup]api.Survey{
	Generic_v1.Group: Generic_v1,
	Wics_v1.Group: Wics_v1,
	Se_soc_v1.Group: Se_soc_v1,
}

var allSurveysByGroupId = map[data.TGroupID]api.Survey{
	"WICS": Wics_v1,
	"SE_SOC": Se_soc_v1,
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

/* GetResponseSimilarity calculates a normalized value of survey response similarity, with 0.0 being the same and 1.0
 * being maximally dissimilar.
 */
func GetResponseSimilarity(group data.SurveyGroup, version int, r1 data.SurveyResponses, r2 data.SurveyResponses) (float64, error) {
	switch {
	case group == Generic_v1.Group && version == Generic_v1.Version:
		return genericV1SimilarityScore(r1, r2), nil
	default:
		return 0.0, errors.Errorf("cannot calculate similarity for survey (%s, %v)", group, version)
	}
}

func OptimalMatching(mentors map[data.TUserID]data.UserSurvey, mentees map[data.TUserID]data.UserSurvey) (mentorMenteeMap map[data.TUserID]data.TUserID, error error) {
	similarity := make(map[data.TUserID]map[data.TUserID]float64)
	for mentorId, mentorSurvey := range mentors {
		for menteeId, menteeSurvey := range mentees {
			if mentorSurvey.Group != menteeSurvey.Group || mentorSurvey.Version != menteeSurvey.Version {
				return nil, errors.Errorf("Survey response type mismatch [(%s, %v) vs (%s, %v)]", mentorSurvey.Group, mentorSurvey.Version, menteeSurvey.Group, menteeSurvey.Version  )
			}
			if sim, err := GetResponseSimilarity(mentorSurvey.Group, mentorSurvey.Version, mentorSurvey.Responses, menteeSurvey.Responses); err != nil {
				return nil, err
			} else {
				similarity[mentorId][menteeId] = sim
			}
		}
	}
	mentorMenteeMap = make(map[data.TUserID]data.TUserID)
	menteeMentorMap := make(map[data.TUserID]data.TUserID)
}
