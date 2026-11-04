package handlers

import (
	"sync"
	"time"
)

type Session struct {
	UserID    int
	ExpiresAt time.Time
}
var (
	sessions   = make(map[string]Session)
	sessionsMu sync.RWMutex
)