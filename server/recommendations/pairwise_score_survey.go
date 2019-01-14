package recommendations

import "letstalk/server/data"

type SurveyPairwiseScore struct {
	Normalize bool
}

func (s *SurveyPairwiseScore) RequiredObjects() []string {
	return []string{"UserSurveys"}
}

func getSurveyMap(user *data.User) map[data.SurveyGroup]data.UserSurvey {
	surveyMap := make(map[data.SurveyGroup]data.UserSurvey)
	if user.UserSurveys != nil {
		for _, survey := range user.UserSurveys {
			surveyMap[survey.Group] = survey
		}
	}
	return surveyMap
}

func (s *SurveyPairwiseScore) Calculate(userOne *data.User, userTwo *data.User) (float32, error) {
	userOneSurveyMap := getSurveyMap(userOne)
	userTwoSurveyMap := getSurveyMap(userTwo)

	var (
		numAnsweredByBoth  uint = 0
		numMatchingAnswers uint = 0
	)

	for group, surveyOne := range userOneSurveyMap {
		if surveyTwo, ok := userTwoSurveyMap[group]; ok {
			for key, responseOne := range surveyOne.Responses {
				if responseTwo, ok := surveyTwo.Responses[key]; ok {
					numAnsweredByBoth++
					if responseOne == responseTwo {
						numMatchingAnswers++
					}
				}
			}
		}
	}

	if numAnsweredByBoth == 0 {
		return 0.0, nil
	} else if s.Normalize {
		return float32(numMatchingAnswers) / float32(numAnsweredByBoth), nil
	} else {
		return float32(numMatchingAnswers), nil
	}
}
