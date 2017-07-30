package main

import (
	"database/sql"
)

const (
	FORMS_FIND_BY_ID_SQL      = `select * from forms where id = $1`
	FORMS_FIND_BY_UUID_SQL    = `select * from forms where uuid = $1`
	FORMS_FIND_BY_USER_ID_SQL = `select * from forms where user_id = $1`
	FORMS_INSERT_SQL          = `insert into forms (user_id, name, description, recaptcha_secret_key) values ($1, $2, $3, $4) returning *`
	FORMS_UPDATE_SQL          = `update forms set name = $2, description = $3, recaptcha_secret_key = $4 where id = $1 returning *`
	FORMS_DELETE_SQL          = `delete from forms where id = $1`
)

type FormsRepository struct {
	db *sql.DB
}

func NewFormsRepository(db *sql.DB) *FormsRepository {
	return &FormsRepository{
		db: db,
	}
}

func (repo *FormsRepository) FindByUserId(userId int) ([]Form, error) {
	var forms []Form

	rows, err := repo.db.Query(FORMS_FIND_BY_USER_ID_SQL, userId)
	for rows.Next() {
		form := new(Form)
		err = repo.scanRow(rows, form)
		forms = append(forms, *form)
	}

	return forms, err
}

func (repo *FormsRepository) Create(userId int, name string, description, recaptchaSecretKey *string) (*Form, error) {
	form := new(Form)
	row := repo.db.QueryRow(FORMS_INSERT_SQL, userId, name, description, recaptchaSecretKey)
	err := repo.scanRow(row, form)
	return form, err
}

func (repo *FormsRepository) FindById(id int) (*Form, error) {
	form := new(Form)
	row := repo.db.QueryRow(FORMS_FIND_BY_ID_SQL, id)
	err := repo.scanRow(row, form)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return form, err
}

func (repo *FormsRepository) FindByUuid(uuid string) (*Form, error) {
	form := new(Form)
	row := repo.db.QueryRow(FORMS_FIND_BY_UUID_SQL, uuid)
	err := repo.scanRow(row, form)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return form, err
}

func (repo *FormsRepository) Update(id int, name string, description, recaptchaSecretKey *string) (*Form, error) {
	form := new(Form)
	row := repo.db.QueryRow(FORMS_UPDATE_SQL, id, name, description, recaptchaSecretKey)
	err := repo.scanRow(row, form)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return form, err
}

func (repo *FormsRepository) Delete(id int) (sql.Result, error) {
	return repo.db.Exec(FORMS_DELETE_SQL, id)
}

type Scannable interface {
	Scan(...interface{}) error
}

func (repo *FormsRepository) scanRow(row Scannable, form *Form) error {
	return row.Scan(&form.id, &form.userId, &form.uuid, &form.name, &form.description, &form.createdAt, &form.updatedAt, &form.recaptchaSecretKey)
}
