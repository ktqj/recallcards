package main

import (
  "net/http"
  "fmt"
  "log"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hi")  
}

func main() {
  http.HandleFunc("/", indexHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}

