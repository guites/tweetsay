package model

import "github.com/dghubble/go-twitter/twitter"

type User struct {
	User *twitter.User
	Active bool
}