package api

type ReqMatchingInfoFlag uint

const (
	REQ_MATCHING_INFO_FLAG_NONE      ReqMatchingInfoFlag = 0
	REQ_MATCHING_INFO_FLAG_AUTH_DATA ReqMatchingInfoFlag = 1 << iota
	REQ_MATCHING_INFO_FLAG_CREDENTIAL
)
