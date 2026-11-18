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
    
    // общие
    http.HandleFunc("/ping", handlers.PingHandler)
    http.HandleFunc("/about", h.About)
    
    // авторизация
    http.HandleFunc("/register", h.RegisterHandler)
    http.HandleFunc("/login", h.LoginHandler)
    http.HandleFunc("/logout", h.LogoutHandler)


    // с учеткой
    http.HandleFunc("/", h.RequireAuth(h.Index))    
    http.HandleFunc("/add", h.RequireAuth(h.AddExpense))
    http.HandleFunc("/expenses/", h.RequireAuth(h.ExpenseRoutes))
   
    http.HandleFunc("/api/stats/by-category", h.RequireAuth(h.APIStatsByCategory))
    http.HandleFunc("/api/stats/by-month", h.RequireAuth(h.APIStatsByMonth))

    http.HandleFunc("/stats", h.RequireAuth(h.StatsPage)) 

    

    
    staticFS := http.FS(static.FS)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticFS)))


    fmt.Println("server is listening on http://localhost:8080 ...")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Printf("server error: %s", err)
    }
}