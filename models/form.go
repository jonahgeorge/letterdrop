package models

import "time"

type Form struct {
	Id                 int
	UserId             int
	Uuid               string
	Name               string
	Description        *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	RecaptchaSecretKey *string
}

func (f *Form) CanCreate(u *User) bool {
	return f.UserId == u.Id
}

func (f *Form) CanView(u *User) bool {
	return f.UserId == u.Id
}

func (f *Form) CanUpdate(u *User) bool {
	return f.UserId == u.Id
}

func (f *Form) CanDelete(u *User) bool {
	return f.UserId == u.Id
}

func (f *Form) FromRow(row Scannable) error {
	return row.Scan(&f.Id, &f.UserId, &f.Uuid, &f.Name, &f.Description, &f.CreatedAt,
		&f.UpdatedAt, &f.RecaptchaSecretKey)
}
