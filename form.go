package main

import (
	"time"
)

type Form struct {
	id          int
	userId      int
	name        string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}
