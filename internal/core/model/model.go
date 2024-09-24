package model

type User struct {
	id           int64
	email        string
	username     string
	passwordHash string
}
