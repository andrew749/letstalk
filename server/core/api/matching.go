package api

type PostMatchingRequest struct {
	Mentor int
	Mentee int
}

type MatchingState int

const (
	MATCHING_STATE_UNVERIFIED MatchingState = iota
	MATCHING_STATE_VERIFIED
	MATCHING_STATE_EXPIRED
)

type MatchingResult struct {
	Mentor int
	Mentee int
	Secret string
	State MatchingState
}
