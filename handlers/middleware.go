package handlers

import (
	"net/http"
	"time"
)

func (h *Handler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		sessionsMu.RLock()
		sess, exists := sessions[cookie.Value]
		sessionsMu.RUnlock()

		if !exists || time.Now().After(sess.ExpiresAt) {
			if exists {
				sessionsMu.Lock()
				delete(sessions, cookie.Value)
				sessionsMu.Unlock()
			}
			
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

func getUserID(r *http.Request) int {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0
	}
	sessionsMu.RLock()
	defer sessionsMu.RUnlock()
	
	sess, exists := sessions[cookie.Value]
	if !exists || time.Now().After(sess.ExpiresAt) {
		return 0
	}
	return sess.UserID
}