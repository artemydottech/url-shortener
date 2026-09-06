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

const reserveAttempts = 5

func reserveCode(url string) (string, error) {
	for attempt := 0; attempt < reserveAttempts; attempt++ {
		code, err := generateCode()
		if err != nil {
			return "", err
		}
		if storage.Save(code, url) {
			return code, nil
		}
	}
	return "", fmt.Errorf("no free code after %d attempts", reserveAttempts)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("POST /api/shorten called")

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

	code, err := reserveCode(req.URL)
	if err != nil {
		log.Printf("Code generation failed: %v", err)
		http.Error(w, "Could not generate a code", http.StatusInternalServerError)
		return
	}

	log.Printf("Created code %s for %s", code, req.URL)

	w.Header().Set("Content-Type", "application/json")

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
