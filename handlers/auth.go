package handlers

import (
	"net/http"
	"project/models"
	"golang.org/x/crypto/bcrypt"
	"time"

	"crypto/rand"
	"encoding/hex"

)

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "GET" {
		pages["register"].ExecuteTemplate(w, "layout", map[string]any{"Title": "Регистрация"})
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		http.Error(w, "Ошибка хэширования", http.StatusInternalServerError)
		return
	}

	err = h.db.CreateUser(models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	})
	if err != nil {
		http.Error(w, "Ошибка базы (возможно, email занят)", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {


	if r.Method == "GET" {
		pages["login"].ExecuteTemplate(w, "layout", map[string]any{"Title": "Вход"})
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	u, err := h.db.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	// делаем куки
	b := make([]byte, 32)
	rand.Read(b)
	sessionID := hex.EncodeToString(b)

	
	sessionsMu.Lock()
	sessions[sessionID] = Session{
		UserID:    int(u.ID),
		ExpiresAt: time.Now().Add(24 * time.Hour), // сессия живет сутки
	}
	sessionsMu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   24 * 60 * 60,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}


func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session_id")
	if err == nil {
	
		sessionsMu.Lock()
		delete(sessions, cookie.Value)
		sessionsMu.Unlock()
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, 
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}