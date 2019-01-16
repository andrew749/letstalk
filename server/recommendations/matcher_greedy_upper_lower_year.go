package recommendations

import (
	"sort"

	"letstalk/server/data"
)

type GreedyUpperLowerYearMatcher struct {
	MaxLowerYears     uint // Maximum number of lower-years per upper-year
	MaxUpperYears     uint // Maximum number of upper-years per lower-year
	YoungestUpperYear uint // Graduating year of "youngest" upper-year
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

func lowerYearsOrderedByDecreasingBestScore(
	lowerYearIterIdxs map[data.TUserID]uint,
	matchMap map[data.TUserID][]UserMatch,
) []data.TUserID {
	matchesToSort := make([]UserMatch, 0)

	for lowerYearId := range lowerYearIterIdxs {
		idx := lowerYearIterIdxs[lowerYearId]
		if idx < uint(len(matchMap[lowerYearId])) {
			matchesToSort = append(matchesToSort, matchMap[lowerYearId][idx])
		}
	}

	sort.Sort(byScore(matchesToSort))

	sortedLowerYearsIds := make([]data.TUserID, len(matchesToSort))
	for i, match := range matchesToSort {
		sortedLowerYearsIds[i] = match.UserOneId
	}

	return sortedLowerYearsIds
}

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
		lowerYearIterIdxs  = make(map[data.TUserID]uint) // keep track of match iter idxs per lower year
		upperYearMatchCnts = make(map[data.TUserID]uint) // keep track of match count per upper year
	)

	for _, lowerYear := range lowerYears {
		lowerYearIterIdxs[lowerYear.UserId] = 0
	}
	for _, upperYear := range upperYears {
		upperYearMatchCnts[upperYear.UserId] = 0
	}

	// Try up to MaxUpperYears matches for each lower year.
	// We go one match at a time for each lower year to make it more fair for lower years that don't
	// have as many "compatible" upper years. Still want to give them decent matches.
	for matchRun := uint(0); matchRun < m.MaxUpperYears; matchRun++ {
		lowerYearIds := lowerYearsOrderedByDecreasingBestScore(lowerYearIterIdxs, matchMap)
		for _, lowerYearId := range lowerYearIds {
			iterIdx := lowerYearIterIdxs[lowerYearId]
			lowerYearMatches := matchMap[lowerYearId]
			for iterIdx < uint(len(lowerYearMatches)) {
				match := lowerYearMatches[iterIdx]
				iterIdx++
				if upperYearMatchCnts[match.UserTwoId] < m.MaxLowerYears {
					// Upper year still has capacity, add this match
					matches = append(matches, match)
					upperYearMatchCnts[match.UserTwoId]++
					break
				}
			}
			lowerYearIterIdxs[lowerYearId] = iterIdx
		}
	}

	return matches, nil
}
