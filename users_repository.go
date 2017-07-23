package main

import (
	"database/sql"
)

const (
	USERS_FIND_BY_ID_SQL                 = "select * from users where id = $1"
	USERS_FIND_BY_EMAIL_AND_PASSWORD_SQL = "select * from users where email = lower($1) and password_digest = crypt($2, password_digest)"
	USERS_INSERT_SQL                     = "insert into users (email, password_digest) values (lower($1), crypt($2, gen_salt('bf', 8))) returning *"
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
	err := repo.db.QueryRow(USERS_FIND_BY_ID_SQL, id).
		Scan(&user.id, &user.email, &user.passwordDigest, &user.createdAt, &user.updatedAt)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (repo *UsersRepository) FindByEmailAndPassword(email, password string) (*User, error) {
	user := new(User)
	err := repo.db.QueryRow(USERS_FIND_BY_EMAIL_AND_PASSWORD_SQL, email, password).
		Scan(&user.id, &user.email, &user.passwordDigest, &user.createdAt, &user.updatedAt)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (repo *UsersRepository) Create(email, password string) (*User, error) {
	user := new(User)
	err := repo.db.QueryRow(USERS_INSERT_SQL, email, password).
		Scan(&user.id, &user.email, &user.passwordDigest, &user.createdAt, &user.updatedAt)
	return user, err
}
