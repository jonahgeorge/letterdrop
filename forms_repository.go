package main

import (
  "database/sql"
  _ "github.com/lib/pq"
)

const (
)

type FormsRepository struct {
  db *sql.DB
}

func NewFormsRepository(db *sql.DB) *FormsRepository {
  return &FormsRepository{
    db: db,
  }
}

func (repo *FormsRepository) FindById(id int) *Form {
  return nil
}
