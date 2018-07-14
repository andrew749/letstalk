package api

import "letstalk/server/data"

type MatchingInfoFlag uint

const (
	MATCHING_INFO_FLAG_NONE      MatchingInfoFlag = 0
	MATCHING_INFO_FLAG_AUTH_DATA MatchingInfoFlag = 1 << iota
	MATCHING_INFO_FLAG_COHORT
)

type Matching struct {
	Mentor data.TUserID       `json:"mentor" binding:"required"`
	Mentee data.TUserID       `json:"mentee" binding:"required"`
	State  data.MatchingState `json:"state"`
}
