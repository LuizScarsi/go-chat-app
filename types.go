package main

import "math/rand"

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	NickName  string `json:"nickName"`
}

func NewUser(firstName, lastName, nickName string) *User {
	return &User{
		ID:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		NickName:  nickName,
	}
}

type Account struct {
	User      *User  `json:"user"`
	CreatedAt string `json:"createdAt"`
}
