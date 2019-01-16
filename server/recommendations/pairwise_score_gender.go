package recommendations

import "letstalk/server/data"

// If both users are either male or female, we try to avoid male-female matches.
// If either didn't specify, we assume that the person is comfortable with any other gender.
type GenderPairwiseScore struct {
}

func (s GenderPairwiseScore) RequiredObjects() []string {
	return []string{}
}

func isGenderMF(user *data.User) bool {
	return user.Gender == data.GENDER_MALE || user.Gender == data.GENDER_FEMALE
}

func (s GenderPairwiseScore) Calculate(userOne *data.User, userTwo *data.User) (float32, error) {
	userOneMF := isGenderMF(userOne)
	userTwoMF := isGenderMF(userTwo)

	// Don't give credit for male-female matches. Others are fine.
	if userOneMF && userTwoMF {
		if userOne.Gender == userTwo.Gender {
			return 1.0, nil
		} else {
			return 0.0, nil
		}
	} else {
		return 1.0, nil
	}
}
