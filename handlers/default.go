package handlers

import (
    "fmt"
    "net/http"
    "html/template"
	"project/models"
	"strconv"
	"time"
	"project/templates"
)

var pages = map[string]*template.Template{
	"index": template.Must(template.ParseFS(templates.FS,"layout.html", "index.html")),
	"about": template.Must(template.ParseFS(templates.FS,"layout.html", "about.html")),
}

func Handler(w http.ResponseWriter, r *http.Request) {
   data := map[string]any{
		"Title":    "МОИ ТРАТЫ",
		"Expenses": models.GetExpenses(),
	}
	pages["index"].ExecuteTemplate(w, "layout", data)
}

func Abouthandler(w http.ResponseWriter, r *http.Request) {
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
        w.WriteHeader(http.StatusMethodNotAllowed) // 405
    }
}
func AddExpenseHandler(w http.ResponseWriter, r *http.Request) {
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

	models.AddExpense(models.Expense{
		Amount:      amount,
		Description: description,
		Date:        date,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}