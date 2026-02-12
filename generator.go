package main

import (
	"crypto/rand"
	"encoding/base64"
)

func generateCode() string {
    bytes := make([]byte, 6) 
    _, err := rand.Read(bytes)
    if err != nil {
        panic(err)  
    }
    code := base64.RawURLEncoding.EncodeToString(bytes)
    return code[:6] 
}
