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

func (s *URLStorage) Save(code, url string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.urls[code] = url
}

func (s *URLStorage) Get(code string) (string, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    url, ok := s.urls[code]
    return url, ok
}
