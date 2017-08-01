package models

import "time"

type Submission struct {
	Id        int
	FormId    int
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Submission) FromRow(row Scannable) error {
	return row.Scan(&s.Id, &s.FormId, &s.Body, &s.CreatedAt, &s.UpdatedAt)
}
