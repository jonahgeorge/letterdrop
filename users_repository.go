package main

import (
	"database/sql"
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
	return &UsersRepository{
		db: db,
	}
}

func (repo *UsersRepository) FindById(id int) (*User, error) {
	user := new(User)
	row := repo.db.QueryRow(USERS_FIND_BY_ID_SQL, id)
	err := repo.scanRow(row, user)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (repo *UsersRepository) FindByEmailAndPassword(email, password string) (*User, error) {
	user := new(User)
	row := repo.db.QueryRow(USERS_FIND_BY_EMAIL_AND_PASSWORD_SQL, email, password)
	err := repo.scanRow(row, user)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (repo *UsersRepository) FindByEmailConfirmationToken(token string) (*User, error) {
	user := new(User)
	row := repo.db.QueryRow(USERS_FIND_BY_EMAIL_CONFIRMATION_TOKEN_SQL, token)
	err := repo.scanRow(row, user)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (repo *UsersRepository) FindByEmail(email string) (*User, error) {
	user := new(User)
	row := repo.db.QueryRow(USERS_FIND_BY_EMAIL_SQL, email)
	err := repo.scanRow(row, user)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (repo *UsersRepository) Create(user *User) (*User, error) {
	row := repo.db.QueryRow(USERS_INSERT_SQL, user.name, user.email, user.passwordDigest, user.emailConfirmationToken)
	err := repo.scanRow(row, user)
	return user, err
}

func (repo *UsersRepository) Update(user *User) (*User, error) {
	row := repo.db.QueryRow(USERS_UPDATE_SQL, user.id, user.name, user.email, user.passwordDigest, user.emailConfirmationToken, user.isEmailConfirmed)
	err := repo.scanRow(row, user)
	return user, err
}

func (repo *UsersRepository) scanRow(row *sql.Row, user *User) error {
	return row.Scan(&user.id, &user.name, &user.email, &user.passwordDigest, &user.createdAt, &user.updatedAt, &user.emailConfirmationToken, &user.isEmailConfirmed)
}
