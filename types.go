package main

import (
	"time"
)

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
	// Email       string `json:"email"`
}

type Account struct {
	AccountID int `json:"accountId"`
	// Login        string    `json:"login"`
	// PasswordHash string    `json:"passwordHash"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	NickName  string `json:"nickName"`
	// AccountType  string    `json:"accountType"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName, nickName string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		NickName:  nickName,
		// AccountType: accountType,
		CreatedAt: time.Now().UTC(),
	}
}
