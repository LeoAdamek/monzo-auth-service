package models

import "github.com/satori/go.uuid"


// Token represents the session token
type Token struct {
	ID string `json:"id"`
	
	Token string `json:"token"`
	StateToken string `json:"-"`
	
	LoginURL string `json:"login_url"`
}

// NewToken generates a new token
func NewToken() Token {
	return Token {
		ID: generateID(),
	}
}

func generateID() string {
	return uuid.NewV4().String()
}