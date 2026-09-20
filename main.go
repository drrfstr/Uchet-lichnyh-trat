package main

import (
    "fmt"
    "net/http"
    "project/handlers" 
)

func main() {
 
    http.HandleFunc("/", handlers.Handler)
    http.HandleFunc("/about", handlers.Abouthandler)
    http.HandleFunc("/ping", handlers.Pinghandler)
    http.HandleFunc("/add", handlers.AddExpenseHandler)

    fmt.Println("server is listening on port 8080...")
    http.ListenAndServe(":8080", nil)
 
}