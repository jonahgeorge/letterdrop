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

func (user *User) FromRow(row Scannable) error {
	return row.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordDigest, &user.CreatedAt,
		&user.UpdatedAt, &user.EmailConfirmationToken, &user.IsEmailConfirmed)
}
