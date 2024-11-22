package main

import (
    "encoding/json"
    "net/http"
)

type Response struct {
    Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    response := Response{Message: "Hello from Netlify!"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    http.HandleFunc("/", handler)
}
