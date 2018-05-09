package api

type MatchingState int

const (
	MATCHING_STATE_UNVERIFIED MatchingState = iota
	MATCHING_STATE_VERIFIED
	MATCHING_STATE_EXPIRED
)

type Matching struct {
	Mentor int `json:"required"`
	Mentee int `json:"required"`
	Secret string
	State MatchingState
}
