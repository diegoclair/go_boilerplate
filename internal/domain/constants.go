package domain

import "time"

// A wait, never a deactivation: the counter expires on its own, so the way back
// in costs the account nothing.
const (
	MaxLoginAttempts   = 5
	LoginAttemptWindow = time.Hour
)
