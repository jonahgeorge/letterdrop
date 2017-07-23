package main

import (
	"database/sql"
	"log"
)

const (
	FORMS_FIND_BY_ID_SQL      = "select * from forms where id = $1"
	FORMS_FIND_BY_USER_ID_SQL = "select * from forms where user_id = $1"
	FORMS_INSERT_SQL          = "insert into forms (user_id, name) values ($1, $2) returning *"
)

type FormsRepository struct {
	db *sql.DB
}

func NewFormsRepository(db *sql.DB) *FormsRepository {
	return &FormsRepository{
		db: db,
	}
}

func (repo *FormsRepository) FindByUserId(userId int) []Form {
	var forms []Form

	rows, err := repo.db.Query(FORMS_FIND_BY_USER_ID_SQL, userId)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		form := new(Form)
		err := rows.Scan(&form.id, &form.userId, &form.name, &form.createdAt, &form.updatedAt)
		if err != nil {
			log.Fatal(err)
		}
		forms = append(forms, *form)
	}

	return forms
}

func (repo *FormsRepository) Create(userId int, name string) (*Form, error) {
	form := new(Form)

	err := repo.db.QueryRow(FORMS_INSERT_SQL, userId, name).Scan(&form.id, &form.userId, &form.name, &form.createdAt, &form.updatedAt)
	if err != nil {
		return nil, err
	}

	return form, nil
}

func (repo *FormsRepository) FindById(id int) *Form {
	form := new(Form)

	err := repo.db.QueryRow(FORMS_FIND_BY_ID_SQL, id).Scan(&form.id, &form.userId, &form.name, &form.createdAt, &form.updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// there were no rows, but otherwise no error occurred
			return nil
		} else {
			log.Fatal(err)
		}
	}

	return form
}
