package repositories

import (
	"database/sql"

	"github.com/jonahgeorge/letterdrop/models"
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

func (repo *SubmissionsRepository) FindByFormId(formId int) ([]models.Submission, error) {
	var submissions []models.Submission
	rows, err := repo.db.Query(SUBMISSIONS_FIND_BY_FORM_ID_SQL, formId)

	for rows.Next() {
		submission := new(models.Submission)
		err = submission.FromRow(rows)
		submissions = append(submissions, *submission)
	}

	return submissions, err
}

func (repo *SubmissionsRepository) Create(formId int, body string) (*models.Submission, error) {
	submission := new(models.Submission)
	row := repo.db.QueryRow(SUBMISSIONS_INSERT_SQL, formId, body)
	err := submission.FromRow(row)
	return submission, err
}

func (repo *SubmissionsRepository) Delete(id int) (sql.Result, error) {
	return repo.db.Exec(SUBMISSIONS_DELETE_SQL, id)
}
