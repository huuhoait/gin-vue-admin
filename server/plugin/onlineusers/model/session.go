package model

import "time"

// SessionInfo is the in-memory representation of a live admin session stored
// in Redis. Token is internal: it powers the kick path (JWT blacklist) and
// must never be returned to API consumers.
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
	Token       string    `json:"-"`
}
