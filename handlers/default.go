package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"project/repository"
	"project/templates"
)

var pages = map[string]*template.Template{
	"index": template.Must(template.ParseFS(templates.FS, "layout.html", "index.html")),
	"about": template.Must(template.ParseFS(templates.FS, "layout.html", "about.html")),
	"edit":  template.Must(template.ParseFS(templates.FS, "layout.html", "edit.html")),
}

type Handler struct {
	db *repository.ExpenseRepository
}

func NewHandler(repo *repository.ExpenseRepository) *Handler {
	return &Handler{db: repo}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.db.GetAll()
	if err != nil {
		http.Error(w, "ошибка базы: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title":    "МОИ ТРАТЫ",
		"Expenses": expenses,
	}
	pages["index"].ExecuteTemplate(w, "layout", data)
}

func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title": "О программе",
	}
	pages["about"].ExecuteTemplate(w, "layout", data)
}

func Pinghandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "pong")
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) AddExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		http.Error(w, "неверная сумма", http.StatusBadRequest)
		return
	}

	description := r.FormValue("description")

	dateStr := r.FormValue("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "неверная дата", http.StatusBadRequest)
		return
	}

	_, err = h.db.Create(amount, description, date)
	if err != nil {
		http.Error(w, "ошибка базы: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) ExpenseRoutes(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case strings.HasSuffix(path, "/edit") && r.Method == "POST":
		h.UpdateExpense(w, r)
	case strings.HasSuffix(path, "/delete") && r.Method == "POST":
		h.DeleteExpense(w, r)
	case r.Method == "GET":
		h.EditForm(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) EditForm(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/expenses/")
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}

	expense, err := h.db.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data := map[string]any{
		"Title":   "Редактирование траты",
		"Expense": expense,
	}
	pages["edit"].ExecuteTemplate(w, "layout", data)
}

func (h *Handler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/expenses/")
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}

	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		http.Error(w, "неверная сумма", http.StatusBadRequest)
		return
	}

	description := r.FormValue("description")

	dateStr := r.FormValue("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "неверная дата", http.StatusBadRequest)
		return
	}

	err = h.db.Update(id, amount, description, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/expenses/")
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}

	err = h.db.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func parseIDFromPath(path, prefix string) (int64, error) {
	raw := strings.TrimPrefix(path, prefix)
	if i := strings.Index(raw, "/"); i != -1 {
		raw = raw[:i]
	}
	return strconv.ParseInt(raw, 10, 64)
}