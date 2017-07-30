package main

import (
	"time"
)

type Form struct {
	id                 int
	userId             int
	uuid               string
	name               string
	description        *string
	createdAt          time.Time
	updatedAt          time.Time
	recaptchaSecretKey *string
}

func (f *Form) CanCreate(u *User) bool {
	return f.userId == u.id
}

func (f *Form) CanView(u *User) bool {
	return f.userId == u.id
}

func (f *Form) CanUpdate(u *User) bool {
	return f.userId == u.id
}

func (f *Form) CanDelete(u *User) bool {
	return f.userId == u.id
}
