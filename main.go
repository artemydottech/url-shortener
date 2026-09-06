package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

var storage = NewStorage()

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

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
		return ":8080"
	}
	if !strings.HasPrefix(port, ":") {
		return ":" + port
	}
	return port
}

// baseURL rebuilds the address the caller reached us on, so the short link
// works behind a proxy and on a real host, not only on localhost.
func baseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if forwarded := r.Header.Get("X-Forwarded-Proto"); forwarded != "" {
		scheme = forwarded
	}
	return scheme + "://" + r.Host
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
		"short_url":  fmt.Sprintf("%s/%s", baseURL(r), code),
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
