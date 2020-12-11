package utils

import (
	"manage/session"
	_ "manage/session/provider"
)

var GlobalSessions *session.Manager

func init() {
	GlobalSessions, _ = session.NewManager("memory", "gosessionid", 7*86400)
	go GlobalSessions.GC()
}
