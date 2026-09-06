package main

import "sync"

type URLStorage struct {
	urls map[string]string
	mu   sync.RWMutex
}

func NewStorage() *URLStorage {
	return &URLStorage{
		urls: make(map[string]string),
	}
}

// Save stores the pair and reports whether the code was free. Overwriting a
// taken code would silently repoint somebody else's existing short link.
func (s *URLStorage) Save(code, url string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, taken := s.urls[code]; taken {
		return false
	}
	s.urls[code] = url
	return true
}

func (s *URLStorage) Get(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.urls[code]
	return url, ok
}
