package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChatMessageRequest struct {
	Message string `json:"message"`
}

type ChatMessage struct {
	NickName string `json:"nickName"`
	Message  string `json:"message"`
}

type CreateAccountRequest struct {
	// AccountType string `json:"accountType"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	NickName  string `json:"nickName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Account struct {
	AccountID    int       `json:"accountId"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"_"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	NickName     string    `json:"nickName"`
	CreatedAt    time.Time `json:"createdAt"`
	// Login        string    `json:"login"`
	// AccountType  string    `json:"accountType"`
}

func (acc *Account) ValidatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(acc.PasswordHash), []byte(pw)) == nil
}

func NewAccount(email, firstName, lastName, nickName, password string) (*Account, error) {
	encryptedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		NickName:     nickName,
		PasswordHash: string(encryptedPw),
		// AccountType: accountType,
		CreatedAt: time.Now().UTC(),
	}, nil
}
