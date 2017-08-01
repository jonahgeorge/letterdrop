package models

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

func (s *Submission) FromRow(row Scannable) error {
	return row.Scan(&s.id, &s.formId, &s.body, &s.createdAt, &s.updatedAt)
}
