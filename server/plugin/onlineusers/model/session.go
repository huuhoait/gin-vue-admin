package model

import "time"

// SessionInfo is the in-memory representation of a live admin session stored
// in Redis. Token is internal — it powers the kick path (JWT blacklist) so it
// MUST be persisted to Redis (which is why the tag is "token", not "-"); the
// API layer wipes it before returning sessions to HTTP callers.
//
// Earlier versions used `json:"-"` here, which excluded Token from json.Marshal
// so the JWT was never written to Redis and Kick silently became a no-op for
// the blacklist path. Don't change this back.
type SessionInfo struct {
	UUID        string    `json:"uuid"`
	UserID      uint      `json:"userID"`
	Username    string    `json:"username"`
	NickName    string    `json:"nickName"`
	AuthorityId uint      `json:"authorityId"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"userAgent"`
	LoginAt     time.Time `json:"loginAt"`
	LastSeenAt  time.Time `json:"lastSeenAt"`
	Token       string    `json:"token,omitempty"`
}
