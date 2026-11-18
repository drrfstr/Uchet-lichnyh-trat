package handlers

import (
	"encoding/json"
	"net/http"
)

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) APIStatsByCategory(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	dateFrom := r.URL.Query().Get("from")
	dateTo := r.URL.Query().Get("to")

	stats, err := h.db.GetStatsByCategory(userID, dateFrom, dateTo)
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

func (h *Handler) APIStatsByMonth(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	dateFrom := r.URL.Query().Get("from")
	dateTo := r.URL.Query().Get("to")

	stats, err := h.db.GetStatsByMonth(userID, dateFrom, dateTo)
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, stats)
}