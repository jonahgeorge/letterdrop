package repositories

import (
	"database/sql"

	"github.com/jonahgeorge/letterdrop/models"
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
	return &FormsRepository{db: db}
}

func (repo *FormsRepository) FindByUserId(userId int) ([]models.Form, error) {
	var forms []models.Form

	rows, err := repo.db.Query(FORMS_FIND_BY_USER_ID_SQL, userId)
	for rows.Next() {
		form := new(models.Form)
		err = form.FromRow(rows)
		forms = append(forms, *form)
	}

	return forms, err
}

func (repo *FormsRepository) Create(userId int, name string, description, recaptchaSecretKey *string) (*models.Form, error) {
	form := new(models.Form)
	row := repo.db.QueryRow(FORMS_INSERT_SQL, userId, name, description, recaptchaSecretKey)
	err := form.FromRow(row)
	return form, err
}

func (repo *FormsRepository) FindById(id int) (*models.Form, error) {
	form := new(models.Form)
	row := repo.db.QueryRow(FORMS_FIND_BY_ID_SQL, id)
	err := form.FromRow(row)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return form, err
}

func (repo *FormsRepository) FindByUuid(uuid string) (*models.Form, error) {
	form := new(models.Form)
	row := repo.db.QueryRow(FORMS_FIND_BY_UUID_SQL, uuid)
	err := form.FromRow(row)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return form, err
}

func (repo *FormsRepository) Update(id int, name string, description, recaptchaSecretKey *string) (*models.Form, error) {
	form := new(models.Form)
	row := repo.db.QueryRow(FORMS_UPDATE_SQL, id, name, description, recaptchaSecretKey)
	err := form.FromRow(row)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return form, err
}

func (repo *FormsRepository) Delete(id int) (sql.Result, error) {
	return repo.db.Exec(FORMS_DELETE_SQL, id)
}
