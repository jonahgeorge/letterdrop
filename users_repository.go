package main

import (
  "log"
  "database/sql"
  _ "github.com/lib/pq"
)

const (
  USERS_FIND_BY_ID_SQL                 = "select * from users where id = $1"
  USERS_FIND_BY_EMAIL_AND_PASSWORD_SQL = "select * from users where email = lower($1) and password_digest = crypt($2, password_digest)"
  USERS_INSERT_SQL                     = "insert into users (name, email, password_digest) values ($1, lower($2), crypt($3, gen_salt('bf', 8))) returning *"
)

type UsersRepository struct {
  db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
  return &UsersRepository{
    db: db,
  }
}

func (repo *UsersRepository) FindById(id int) *User {
  user := new(User)

  err := repo.db.QueryRow(USERS_FIND_BY_ID_SQL, id).
  Scan(&user.id, &user.email, &user.passwordDigest, &user.createdAt, &user.updatedAt)

  if err != nil {
    if err == sql.ErrNoRows {
      // there were no rows, but otherwise no error occurred
      return nil
    } else {
      log.Fatal(err)
    }
  }

  return user
}

func (repo *UsersRepository) FindByEmailAndPassword(email, password string) *User {
  user := new(User)

  err := repo.db.QueryRow(USERS_FIND_BY_EMAIL_AND_PASSWORD_SQL, email, password).
  Scan(&user.id, &user.email, &user.passwordDigest, &user.createdAt, &user.updatedAt)

  if err != nil {
    if err == sql.ErrNoRows {
      // there were no rows, but otherwise no error occurred
      return nil
    } else {
      log.Fatal(err)
    }
  }

  return user
}

