package main
 
import (
    "fmt"
    "log"
    "net/http"
)
 
func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Привет! Это HTTP-сервер на Go.")
}
 
func aboutHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
 
    html := `
<!DOCTYPE html>
<html lang="ru">
  <head>
    <meta charset="utf-8">
    <title>О нас</title>
  </head>
  <body>
    <h1>О нашей компании</h1>
    <p>Мы разрабатываем сервисы на Go и делаем их простыми, быстрыми и надёжными.</p>
    <a href="/">Вернуться на главную</a> |
    <a href="/ping">Ping</a>
  </body>
</html>
    `
    fmt.Fprintf(w, html)
}
 
func pingHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
 
    html := `
<!DOCTYPE html>
<html lang="ru">
  <head>
    <meta charset="utf-8">
    <title>Ping</title>
  </head>
  <body>
    <h1>Ping</h1>
    <p>Сервер работает: OK.</p>
    <p>Время запроса: ` + r.Header.Get("Date") + `</p>
    <a href="/">Вернуться на главную</a> |
    <a href="/about">О нас</a>
  </body>
</html>
    `
    fmt.Fprintf(w, html)
}
 
func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/about", aboutHandler)
    http.HandleFunc("/ping", pingHandler)
 
    port := ":8080"
    log.Printf("Сервер запущен на http://localhost%s\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}
