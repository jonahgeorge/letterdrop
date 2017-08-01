package models

type Permissible interface {
	CanCreate(*User) bool
	CanView(*User) bool
	CanUpdate(*User) bool
	CanDelete(*User) bool
}
