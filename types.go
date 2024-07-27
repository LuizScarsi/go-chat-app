package main

import "math/rand"

type User struct {
	ID int
	FirstName string
	LastName string
	NickName string
}

func NewUser(firstName, lastName, nickName string) *User {
	return &User{
		ID: rand.Intn(10000),
		FirstName: firstName,
		LastName: lastName,
		NickName: nickName,
	}
}