package main

import (
	"time"
)

type User struct {
	id                     int
	name                   string
	email                  string
	passwordDigest         string
	createdAt              time.Time
	updatedAt              time.Time
	emailConfirmationToken *string
	isEmailConfirmed       bool
}

type Permissible interface {
	CanCreate(*User) bool
	CanView(*User) bool
	CanUpdate(*User) bool
	CanDelete(*User) bool
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
