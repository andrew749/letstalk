package api

import "letstalk/server/data"

type Matching struct {
	Mentor int `json:"mentor" binding:"required"`
	Mentee int `json:"mentee" binding:"required"`
	State data.MatchingState `json:"state"`
}
