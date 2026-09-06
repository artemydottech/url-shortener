package main

import (
	"crypto/rand"
	"encoding/base64"
)

const codeLength = 6

func generateCode() (string, error) {
	bytes := make([]byte, codeLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:codeLength], nil
}
