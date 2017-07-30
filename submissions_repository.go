package main

import (
	"database/sql"
)

const (
	SUBMISSIONS_FIND_BY_FORM_ID_SQL = `select * from submissions where form_id = $1`
	SUBMISSIONS_INSERT_SQL          = `insert into submissions (form_id, body) values ($1, $2) returning *`
	SUBMISSIONS_DELETE_SQL          = `delete from submissions where id = $1`
)

type SubmissionsRepository struct {
	db *sql.DB
}

func NewSubmissionsRepository(db *sql.DB) *SubmissionsRepository {
	return &SubmissionsRepository{
		db: db,
	}
}

func (repo *SubmissionsRepository) FindByFormId(formId int) ([]Submission, error) {
	var submissions []Submission
	rows, err := repo.db.Query(SUBMISSIONS_FIND_BY_FORM_ID_SQL, formId)

	for rows.Next() {
		submission := new(Submission)
		err = repo.scanRow(rows, submission)
		submissions = append(submissions, *submission)
	}

	return submissions, err
}

func (repo *SubmissionsRepository) Create(formId int, body string) (*Submission, error) {
	submission := new(Submission)
	row := repo.db.QueryRow(SUBMISSIONS_INSERT_SQL, formId, body)
	err := repo.scanRow(row, submission)
	return submission, err
}

func (repo *SubmissionsRepository) Delete(id int) (sql.Result, error) {
	return repo.db.Exec(SUBMISSIONS_DELETE_SQL, id)
}

func (repo *SubmissionsRepository) scanRow(row Scannable, submission *Submission) error {
	return row.Scan(&submission.id, &submission.formId, &submission.body, &submission.createdAt, &submission.updatedAt)
}
