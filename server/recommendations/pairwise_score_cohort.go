package recommendations

import "letstalk/server/data"

type CohortPairwiseScore struct {
}

func (s *CohortPairwiseScore) RequiredObjects() []string {
	return []string{"Cohort.Cohort"}
}

func getCohort(user *data.User) *data.Cohort {
	if user.Cohort != nil && user.Cohort.Cohort != nil {
		return user.Cohort.Cohort
	}
	return nil
}

func (s *CohortPairwiseScore) Calculate(userOne *data.User, userTwo *data.User) (float32, error) {
	cohortOne := getCohort(userOne)
	cohortTwo := getCohort(userTwo)

	if cohortOne != nil && cohortTwo != nil && cohortOne.CohortId == cohortTwo.CohortId {
		return 1.0, nil
	} else {
		return 0.0, nil
	}
}
