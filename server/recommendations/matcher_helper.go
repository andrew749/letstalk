package recommendations

import "letstalk/server/data"

// Calculates all user matches by taking the product of two lists of users.
// Will return m * n user matches given m users on the left and n users on the right.
// Returns a list of all of these matches. userOneId will be users from the left list and userTwoId
// will be users from the right list. Matches are unsorted
func calculateSplitUserMatches(
	usersLeft []data.User,
	usersRight []data.User,
	score PairwiseScore,
) ([]UserMatch, error) {
	matches := make([]UserMatch, 0)
	for _, userLeft := range usersLeft {
		for _, userRight := range usersRight {
			// Don't create matches for two of the same user
			if userLeft.UserId != userRight.UserId {
				value, err := score.Calculate(&userLeft, &userRight)
				if err != nil {
					return nil, err
				}
				userMatch := UserMatch{
					UserOneId: userLeft.UserId,
					UserTwoId: userRight.UserId,
					Score:     value,
				}
				matches = append(matches, userMatch)
			}
		}
	}
	return matches, nil
}
