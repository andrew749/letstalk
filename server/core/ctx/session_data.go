package ctx

import (
	"time"
)

/**
 * Stores data related to a certain session.
 */
type SessionData struct {
	SessionId  *string
	ExpiryDate time.Time
}
