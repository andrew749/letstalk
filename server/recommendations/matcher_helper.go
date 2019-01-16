package recommendations

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"

	"letstalk/server/data"
)

type byScore []UserMatch

func (a byScore) Len() int {
	return len(a)
}

func (a byScore) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byScore) Less(i, j int) bool {
	// Sort by decreasing
	return a[i].Score > a[j].Score
}

// Calculates all user matches by taking the product of two lists of users.
// Will return m * n user matches given m users on the left and n users on the right.
// Map is keyed by user ids of left users, and each list is sorted by decreasing score.
func calculateSplitUserMatches(
	usersLeft []data.User,
	usersRight []data.User,
	score PairwiseScore,
) (map[data.TUserID][]UserMatch, error) {
	matches := make(map[data.TUserID][]UserMatch)
	for _, userLeft := range usersLeft {
		matches[userLeft.UserId] = make([]UserMatch, 0, len(usersRight))
		for _, userRight := range usersRight {
			if userLeft.UserId == userRight.UserId {
				return nil, errors.New(
					fmt.Sprintf("User %d is repeated in left and right lists", userLeft.UserId))
			}
			value, err := score.Calculate(&userLeft, &userRight)
			if err != nil {
				return nil, err
			}
			userMatch := UserMatch{
				UserOneId: userLeft.UserId,
				UserTwoId: userRight.UserId,
				Score:     value,
			}
			matches[userLeft.UserId] = append(matches[userLeft.UserId], userMatch)
		}
		// Sort matches by decreasing score
		sort.Sort(byScore(matches[userLeft.UserId]))
	}
	return matches, nil
}
