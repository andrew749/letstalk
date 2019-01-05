package api

type UserState string

/**
 * These states will likely change.
 * Currently a later state implies that the previous states are satisfied
 * This is currently a linear state hierarchy
 */
const (
	ACCOUNT_CREATED        UserState = "account_created"        // first state
	ACCOUNT_EMAIL_VERIFIED UserState = "account_email_verified" // UW email has been verified
	ACCOUNT_HAS_BASIC_INFO UserState = "account_has_basic_info" // the account now has basic user info like cohort
	ACCOUNT_SETUP          UserState = "account_setup"          // the account has enough information to proceed
	ACCOUNT_MATCHED        UserState = "account_matched"        // account has been matched a peer
)
