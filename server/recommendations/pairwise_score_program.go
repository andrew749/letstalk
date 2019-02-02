package recommendations

import "letstalk/server/data"

type ProgramPairwiseScore struct {
}

func (s ProgramPairwiseScore) RequiredObjects() []string {
	return []string{"Cohort.Cohort"}
}

func getProgramId(user *data.User) *string {
	if user.Cohort != nil && user.Cohort.Cohort != nil {
		return &user.Cohort.Cohort.ProgramId
	}
	return nil
}

func (s ProgramPairwiseScore) Calculate(userOne *data.User, userTwo *data.User) (Score, error) {
	programOneId := getProgramId(userOne)
	programTwoId := getProgramId(userTwo)

	if programOneId != nil && programTwoId != nil && *programOneId == *programTwoId {
		return 1.0, nil
	} else {
		return 0.0, nil
	}
}
