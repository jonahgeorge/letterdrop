package repositories

import (
	"database/sql"

	"github.com/jonahgeorge/letterdrop/models"
)

const (
	USERS_FIND_BY_ID_SQL                       = `select * from users where id = $1`
	USERS_FIND_BY_EMAIL_SQL                    = `select * from users where email = lower($1)`
	USERS_FIND_BY_EMAIL_AND_PASSWORD_SQL       = `select * from users where email = lower($1) and password_digest = crypt($2, password_digest)`
	USERS_FIND_BY_EMAIL_CONFIRMATION_TOKEN_SQL = `select * from users where email_confirmation_token = $1`

	USERS_INSERT_SQL = `insert into users (name, email, password_digest, email_confirmation_token) values ($1, lower($2), crypt($3, gen_salt('bf', 8)), $4) returning *`
	USERS_UPDATE_SQL = `update users set name = $2, email = $3, password_digest = $4, updated_at = now(), email_confirmation_token = $5, is_email_confirmed = $6 where id = $1 returning *`
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (repo *UsersRepository) FindById(id int) (*models.User, error) {
	user := new(models.User)
	row := repo.db.QueryRow(USERS_FIND_BY_ID_SQL, id)
	err := user.FromRow(row)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (repo *UsersRepository) FindByEmailAndPassword(email, password string) (*models.User, error) {
	user := new(models.User)
	row := repo.db.QueryRow(USERS_FIND_BY_EMAIL_AND_PASSWORD_SQL, email, password)
	err := user.FromRow(row)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (repo *UsersRepository) FindByEmailConfirmationToken(token string) (*models.User, error) {
	user := new(models.User)
	row := repo.db.QueryRow(USERS_FIND_BY_EMAIL_CONFIRMATION_TOKEN_SQL, token)
	err := user.FromRow(row)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (repo *UsersRepository) FindByEmail(email string) (*models.User, error) {
	user := new(models.User)
	row := repo.db.QueryRow(USERS_FIND_BY_EMAIL_SQL, email)
	err := user.FromRow(row)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (repo *UsersRepository) Create(user *models.User) (*models.User, error) {
	row := repo.db.QueryRow(USERS_INSERT_SQL, user.Name, user.Email,
		user.PasswordDigest, user.EmailConfirmationToken)
	err := user.FromRow(row)
	return user, err
}

func (repo *UsersRepository) Update(user *models.User) (*models.User, error) {
	row := repo.db.QueryRow(USERS_UPDATE_SQL, user.Id, user.Name, user.Email,
		user.PasswordDigest, user.EmailConfirmationToken, user.IsEmailConfirmed)
	err := user.FromRow(row)
	return user, err
}
