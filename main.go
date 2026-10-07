package main

import (
    "fmt"
    "net/http"

    "project/handlers"
    "project/repository"
    "project/static"
)

func main() {
    repo := repository.NewExpenseRepository("expenses.db")
    defer repo.Close()

    h := handlers.NewHandler(repo)

    http.HandleFunc("/", h.Index)
    http.HandleFunc("/about", h.About)
    http.HandleFunc("/ping", handlers.Pinghandler)
    http.HandleFunc("/add", h.AddExpense)
    http.HandleFunc("/expenses/", h.ExpenseRoutes)

    staticFS := http.FS(static.FS)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticFS)))

    fmt.Println("server is listening on http://localhost:8080 ...")
    http.ListenAndServe(":8080", nil)
}