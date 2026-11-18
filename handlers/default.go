package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"project/models"
	"project/repository"
	"project/templates"
)

var pages = map[string]*template.Template{
	"index":    template.Must(template.ParseFS(templates.FS, "layout.html", "index.html")),
	"about":    template.Must(template.ParseFS(templates.FS, "layout.html", "about.html")),
	"edit":     template.Must(template.ParseFS(templates.FS, "layout.html", "edit.html")),
	"register": template.Must(template.ParseFS(templates.FS, "layout.html", "register.html")),
	"login":    template.Must(template.ParseFS(templates.FS, "layout.html", "login.html")),
	"stats": 	template.Must(template.ParseFS(templates.FS, "layout.html", "stats.html")),
}

type Handler struct {
	db *repository.ExpenseRepository
}

func NewHandler(repo *repository.ExpenseRepository) *Handler {
	return &Handler{db: repo}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	userID := getUserID(r)

	user, err := h.db.GetUserByID(userID)
	if err != nil {
		http.Error(w, "ошибка получения пользователя", http.StatusInternalServerError)
		return
	}

	categoryStr := r.URL.Query().Get("category")
	dateFrom := r.URL.Query().Get("from")
	dateTo := r.URL.Query().Get("to")

	var categoryID int = 0
	if categoryStr != "" {
		if cid, err := strconv.Atoi(categoryStr); err == nil {
			categoryID = cid
		}
	}

	expenses, err := h.db.GetAllFiltered(categoryID, dateFrom, dateTo, userID)
	if err != nil {
		http.Error(w, "ошибка базы: "+err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := h.db.GetAllCategories()
	if err != nil {
		http.Error(w, "ошибка базы: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title":           "МОИ ТРАТЫ",
		"Expenses":        expenses,
		"Categories":      categories,
		"CurrentCategory": categoryID,
		"CurrentFrom":     dateFrom,
		"CurrentTo":       dateTo,
		"User": 		   user,
	}

	err = pages["index"].ExecuteTemplate(w, "layout", data)
	if err != nil {
		fmt.Printf("ExecuteTemplate error: %s", err)
	}
}

func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title": "О программе",
	}
	err := pages["about"].ExecuteTemplate(w, "layout", data)
	if err != nil {
		fmt.Printf("ExecuteTemplate error: %s", err)
	}
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "pong")
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) AddExpense(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	
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

	category, err := strconv.Atoi(r.FormValue("category_id"))
	if err != nil {
		http.Error(w, "неверный id категории", http.StatusBadRequest)
	}

	newExpense := models.Expense{
		Amount:      amount,
		Description: description,
		Date:        date,
		CategoryID:  category,
		UserID: userID,
	}

	_, err = h.db.CreateExpense(newExpense)
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
	
	userID := getUserID(r)

	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}

	expense, err := h.db.GetByID(id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	categories, err := h.db.GetAllCategories()
	if err != nil {
		http.Error(w, "ошибка базы: "+err.Error(), http.StatusNotFound)
		return
	}

	data := map[string]any{
		"Title":      "Редактирование траты",
		"Expense":    expense,
		"Categories": categories,
	}
	err = pages["edit"].ExecuteTemplate(w, "layout", data)
	if err != nil {
		fmt.Printf("ExecuteTemplate error: %s", err)
	}
}

func (h *Handler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/expenses/")
	userID := getUserID(r)

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

	category, err := strconv.Atoi(r.FormValue("category_id"))
	if err != nil {
		http.Error(w, "неверная категория", http.StatusBadRequest)
		return
	}

	upd := models.Expense{
		ID:          id,
		Amount:      amount,
		Description: description,
		Date:        date,
		CategoryID:  category,
		UserID: userID,
	}

	err = h.db.Update(upd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/expenses/")
	userID := getUserID(r)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}

	err = h.db.Delete(id, userID)
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

func (h *Handler) StatsPage(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	user, err := h.db.GetUserByID(userID)
	if err != nil {
		http.Error(w, "ошибка получения пользователя", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title": "Статистика",
		"User":  user,
	}
	pages["stats"].ExecuteTemplate(w, "layout", data)
}