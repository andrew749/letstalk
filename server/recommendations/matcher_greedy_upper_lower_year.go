package recommendations

import (
	"container/heap"

	"letstalk/server/data"
)

type UserMatchesPQ []UserMatch

func (pq UserMatchesPQ) Len() int {
	return len(pq)
}

func (pq UserMatchesPQ) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Score > pq[j].Score
}

func (pq UserMatchesPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *UserMatchesPQ) Push(x interface{}) {
	*pq = append(*pq, x.(UserMatch))
}

func (pq *UserMatchesPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	match := old[n-1]
	*pq = old[0 : n-1]
	return match
}

// Returns a "set" of all user ids from given list of users.
func getUserIdSet(users []data.User) map[data.TUserID]interface{} {
	userIdSet := make(map[data.TUserID]interface{})
	for _, user := range users {
		userIdSet[user.UserId] = nil
	}
	return userIdSet
}

type GreedyUpperLowerYearMatcher struct {
	MaxLowerYearsPerUpperYear uint // Maximum number of lower-years per upper-year
	MaxUpperYearsPerLowerYear uint // Maximum number of upper-years per lower-year
	YoungestUpperYear         uint // Graduating year of "youngest" upper-year
}

func (m GreedyUpperLowerYearMatcher) RequiredObjects() []string {
	return []string{"Cohort.Cohort"}
}

// Splits users into upper-year and lower-year
func (m GreedyUpperLowerYearMatcher) splitUsersByGradYear(
	users []data.User,
) ([]data.User, []data.User) {
	lowerYears := make([]data.User, 0)
	upperYears := make([]data.User, 0)

	for _, user := range users {
		if user.Cohort != nil && user.Cohort.Cohort != nil {
			if user.Cohort.Cohort.GradYear > m.YoungestUpperYear {
				lowerYears = append(lowerYears, user)
			} else {
				upperYears = append(upperYears, user)
			}
		}
	}

	return lowerYears, upperYears
}

// Matches users in lower year cohorts to users in upper year cohorts greedily
// We essentially choose the highest scoring matches first.
// Some amount of fairness is maintained by the fact that we make sure that each lower year has
// a match before we give any lower year their next match. Thus, the matching is done in epochs.
func (m GreedyUpperLowerYearMatcher) Match(
	users []data.User,
	score PairwiseScore,
) ([]UserMatch, error) {
	lowerYears, upperYears := m.splitUsersByGradYear(users)
	matchMap, err := calculateSplitUserMatches(lowerYears, upperYears, score)
	if err != nil {
		return nil, err
	}

	var (
		matches            = make([]UserMatch, 0)        // list containing matches
		upperYearMatchCnts = make(map[data.TUserID]uint) // number of matches made per upper year
		matchCandidateHeap = make(UserMatchesPQ, 0)      // pq of matches by decreasing score
	)

	// Init heap with all matches
	for _, match := range matchMap {
		matchCandidateHeap = append(matchCandidateHeap, match)
	}
	heap.Init(&matchCandidateHeap)

	// Init all upper year match counts to 0
	for _, upperYear := range upperYears {
		upperYearMatchCnts[upperYear.UserId] = 0
	}

	// Each "epoch" attempts to find 1 match for each lower year, attempting to perserve some level
	// of "fairness" during matching.
	for matchRun := uint(0); matchRun < m.MaxUpperYearsPerLowerYear; matchRun++ {
		lowerYearsToMatch := getUserIdSet(lowerYears)
		skippedMatches := make([]UserMatch, 0)

		for len(lowerYearsToMatch) > 0 && len(matchCandidateHeap) > 0 {
			match := heap.Pop(&matchCandidateHeap).(UserMatch)
			_, lowerYearNotMatched := lowerYearsToMatch[match.UserOneId]
			upperYearHasCapacity := upperYearMatchCnts[match.UserTwoId] < m.MaxLowerYearsPerUpperYear

			if upperYearHasCapacity {
				if lowerYearNotMatched {
					delete(lowerYearsToMatch, match.UserOneId)
					matches = append(matches, match)
					upperYearMatchCnts[match.UserTwoId]++
				} else {
					// Lower year is already matched for this round, but upper year still has capacity, so
					// we want to place it back on heap for next round.
					skippedMatches = append(skippedMatches, match)
				}
			}
		}

		// Place any skipped matches back on heap
		for _, match := range skippedMatches {
			heap.Push(&matchCandidateHeap, match)
		}
	}

	return matches, nil
}
