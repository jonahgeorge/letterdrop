package main

import (
	"time"
)

type Submission struct {
	id        int
	formId    int
	body      string
	createdAt time.Time
	updatedAt time.Time
}
