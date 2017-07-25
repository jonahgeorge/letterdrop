package main

import (
	"database/sql"
	"log"
)

const (
	FORMS_FIND_BY_ID_SQL = `
	select id, user_id, uuid, name, description, created_at, updated_at 
	from forms 
	where id = $1`

	FORMS_FIND_BY_UUID_SQL = `
	select id, user_id, uuid, name, description, created_at, updated_at 
	from forms 
	where uuid = $1`

	FORMS_FIND_BY_USER_ID_SQL = `
	select id, user_id, uuid, name, description, created_at, updated_at 
	from forms 
	where user_id = $1`

	FORMS_INSERT_SQL = `
	insert into forms (user_id, name, description) 
	values ($1, $2, $3) 
	returning id, user_id, name, description, created_at, updated_at`
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
		err := rows.Scan(&form.id, &form.userId, &form.uuid, &form.name, &form.description, &form.createdAt, &form.updatedAt)
		if err != nil {
			log.Fatal(err)
		}
		forms = append(forms, *form)
	}

	return forms
}

func (repo *FormsRepository) Create(userId int, name, description string) (*Form, error) {
	form := new(Form)
	err := repo.db.QueryRow(FORMS_INSERT_SQL, userId, name, description).
		Scan(&form.id, &form.userId, &form.uuid, &form.name, &form.description, &form.createdAt, &form.updatedAt)
	return form, err
}

func (repo *FormsRepository) FindById(id int) (*Form, error) {
	form := new(Form)
	err := repo.db.QueryRow(FORMS_FIND_BY_ID_SQL, id).
		Scan(&form.id, &form.userId, &form.uuid, &form.name, &form.description, &form.createdAt, &form.updatedAt)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}
	return form, err
}

func (repo *FormsRepository) FindByUuid(uuid string) (*Form, error) {
	form := new(Form)
	err := repo.db.QueryRow(FORMS_FIND_BY_UUID_SQL, uuid).
		Scan(&form.id, &form.userId, &form.uuid, &form.name, &form.description, &form.createdAt, &form.updatedAt)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}
	return form, err
}
