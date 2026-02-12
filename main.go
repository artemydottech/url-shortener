package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type ShortenRequest struct {
    URL string `json:"url"`
}

var storage = NewStorage()

func main() {
    port := getPort()
    
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })
    http.HandleFunc("/api/shorten", shortenHandler)
    http.HandleFunc("/", redirectHandler)
    
    fmt.Printf("(!) Сервер запущен на порту %s\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}

func getPort() string {
    port := os.Getenv("PORT")
    if port == "" {
        port = ":8080"
    }
    return port
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("POST /api/shorten called") 
    
    if r.Method != "POST" {
        http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
        return
    }

    var req ShortenRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("JSON error: %v", err)
        http.Error(w, "Bad JSON", http.StatusBadRequest)
        return
    }

    if req.URL == "" {
        http.Error(w, "URL required", http.StatusBadRequest)
        return
    }

    code := generateCode()
    storage.Save(code, req.URL)
    
    log.Printf("Created code %s for %s", code, req.URL)

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    
    response := map[string]string{
        "short_code": code,
        "short_url":  fmt.Sprintf("http://localhost%s/%s", getPort(), code),
    }
    json.NewEncoder(w).Encode(response)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Path[1:]
    if code == "" {
        http.Error(w, "Code required", http.StatusBadRequest)
        return
    }
    
    url, exists := storage.Get(code)
    if !exists {
        http.Error(w, "Not found", http.StatusNotFound)
        return
    }
    
    http.Redirect(w, r, url, http.StatusMovedPermanently)
}
