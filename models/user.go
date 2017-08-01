package models

import "time"

type User struct {
	Id                     int
	Name                   string
	Email                  string
	PasswordDigest         string
	CreatedAt              time.Time
	UpdatedAt              time.Time
	EmailConfirmationToken *string
	IsEmailConfirmed       bool
}

func (u *User) CanCreate(resource Permissible) bool {
	return resource.CanCreate(u)
}

func (u *User) CanView(resource Permissible) bool {
	return resource.CanView(u)
}

func (u *User) CanUpdate(resource Permissible) bool {
	return resource.CanUpdate(u)
}

func (u *User) CanDelete(resource Permissible) bool {
	return resource.CanDelete(u)
}

func (u *User) FromRow(row Scannable) error {
	return row.Scan(&u.Id, &u.Name, &u.Email, &u.PasswordDigest, &u.CreatedAt,
		&u.UpdatedAt, &u.EmailConfirmationToken, &u.IsEmailConfirmed)
}
