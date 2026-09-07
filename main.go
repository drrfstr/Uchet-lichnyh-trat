package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Привет, это Go-сервер!")
}

func abouthandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Это программа, которая ведет учет личных трат")
}

func pinghandler(w http.ResponseWriter, r *http.Request) {
 
    if r.Method == "GET" {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "pong")

    } else {
        w.WriteHeader(http.StatusMethodNotAllowed) // 405
    }
}

func main() {
 
    http.HandleFunc("/", handler)
    http.HandleFunc("/about", abouthandler)
    http.HandleFunc("/ping", pinghandler)

    fmt.Println("server is listening on port 8080...")
    http.ListenAndServe(":8080", nil)
 
}