package session

import "go.dokimi.dev/assert/prop"

// KeepsTheLastPut checks that Get returns the session of the last Put of its
// ID.
//
//dokimi:check Sessions.Get defect=keepsTheFirstPut
func KeepsTheLastPut(h *SessionsHandle, c *prop.Case) {}

// keepsTheFirstPut returns a Sessions whose Put keeps the first session of
// each ID.
func keepsTheFirstPut(*SessionsHandle) Sessions { return nil }
