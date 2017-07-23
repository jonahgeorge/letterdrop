package main

import (
	"database/sql"
	"log"
)

const (
	FIND_INSERT_SQL     = "insert into submissions (form_id, body) values ($1, $2)"
	FIND_BY_FORM_ID_SQL = "select * from submissions where form_id = $1"
)

type SubmissionsRepository struct {
	db *sql.DB
}

func NewSubmissionsRepository(db *sql.DB) *SubmissionsRepository {
	return &SubmissionsRepository{
		db: db,
	}
}

func (repo *SubmissionsRepository) FindByFormId(formId int) []Submission {
	var submissions []Submission

	rows, err := repo.db.Query(FIND_BY_FORM_ID_SQL, formId)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		submission := new(Submission)
		err := rows.Scan(&submission.id, &submission.formId, &submission.body, &submission.createdAt, &submission.updatedAt)
		if err != nil {
			log.Fatal(err)
		}
		submissions = append(submissions, *submission)
	}

	return submissions
}
