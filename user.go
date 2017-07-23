package main

import (
	"time"
)

type User struct {
	id             int
	email          string
	passwordDigest string
	createdAt      time.Time
	updatedAt      time.Time
}
