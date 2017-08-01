package models

type Scannable interface {
	Scan(...interface{}) error
}
