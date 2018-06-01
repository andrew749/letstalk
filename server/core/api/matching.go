package api

import "letstalk/server/data"

type MatchingInfoFlag uint

const (
	MATCHING_INFO_FLAG_AUTH_DATA MatchingInfoFlag = 1 << iota
	MATCHING_INFO_FLAG_COHORT
)

type Matching struct {
	Mentor int                `json:"mentor" binding:"required"`
	Mentee int                `json:"mentee" binding:"required"`
	State  data.MatchingState `json:"state"`
}
