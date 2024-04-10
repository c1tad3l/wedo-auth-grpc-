package models

type User struct {
	Uuid        string
	Email       string
	PassHash    []byte
	Phone       string
	DateOfBirth string
	Username    string
}
